package influxdb

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/influxdata/influxdb1-client" // needed due to go mod bug
	influxdbv2 "github.com/influxdata/influxdb1-client/v2"

	"github.com/slok/grafterm/internal/model"
	"github.com/slok/grafterm/internal/service/metric"
)

// ConfigGatherer is the configuration of the InfluxDB gatherer.
type ConfigGatherer struct {
	Addr     string
	Database string
	Client   influxdbv2.Client
}

func (c *ConfigGatherer) defaults() error {
	var err error

	if c.Database == "" {
		return fmt.Errorf("no influxdb database given")
	}

	if c.Client == nil {
		c.Client, err = influxdbv2.NewHTTPClient(
			influxdbv2.HTTPConfig{
				Addr: c.Addr,
			})
	}
	return err
}

type gatherer struct {
	cli influxdbv2.Client
	cfg ConfigGatherer
}

// NewGatherer returns a new metric gatherer for influxdb backends.
func NewGatherer(cfg ConfigGatherer) (metric.Gatherer, error) {
	err := cfg.defaults()
	if err != nil {
		return &gatherer{}, err
	}

	return &gatherer{
		cli: cfg.Client,
		cfg: cfg,
	}, nil
}

func (g *gatherer) GatherSingle(ctx context.Context, query model.Query, t time.Time) ([]model.MetricSeries, error) {
	res, err := g.GatherRange(ctx, query, t.Add(-1*5), t, 0)
	if err != nil {
		return []model.MetricSeries{}, err
	}

	if len(res) < 1 {
		return []model.MetricSeries{}, fmt.Errorf("server didn't return any metric series")
	}
	if len(res) > 1 {
		return []model.MetricSeries{}, fmt.Errorf("server returned more than one metric series, got %d", len(res))
	}

	// Get the latest datapoint
	res[0].Metrics = res[0].Metrics[len(res[0].Metrics)-1:]

	return res, nil
}

func (g *gatherer) GatherRange(ctx context.Context, query model.Query, start, end time.Time, _ time.Duration) ([]model.MetricSeries, error) {
	res := []model.MetricSeries{}

	// Get the data from the InfluxDB API
	q := influxdbv2.NewQuery(query.Expr, g.cfg.Database, "ms")
	resp, err := g.cli.Query(q)
	if err != nil {
		return res, err
	}
	if resp.Error() != nil {
		return res, resp.Error()
	}

	// Build the metric series
	for _, result := range resp.Results {
		for _, serie := range result.Series {
			metrics := []model.Metric{}
			for _, value := range serie.Values {
				v, err := value[1].(json.Number).Float64()
				if err != nil {
					return res, err
				}
				t := time.Time{}
				switch value[0].(type) {
				case string:
					t, err = time.Parse("2006-01-02T15:04:05Z", value[0].(string))
					if err != nil {
						return res, err
					}
				case json.Number:
					c, err := value[0].(json.Number).Int64()
					if err != nil {
						return res, err
					}
					t = time.Unix(c/1000, c%1000)
				default:
				}
				m := model.Metric{
					TS:    t,
					Value: v,
				}
				metrics = append(metrics, m)
			}
			label := serie.Name
			//TODO(rochaporto): Allow alias based on a tag value
			//if serie.Tags != nil {
			//	if v1, ok := serie.Tags["version"]; ok {
			//		label = v1
			//	}
			//}
			res = append(res, model.MetricSeries{ID: label, Metrics: metrics})
		}
	}

	return res, nil
}
