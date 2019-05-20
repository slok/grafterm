package graphite

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	graphite "github.com/JensRantil/graphite-client"

	"github.com/slok/grafterm/internal/model"
	"github.com/slok/grafterm/internal/service/metric"
)

// ConfigGatherer is the configuration of the Graphite gatherer.
type ConfigGatherer struct {
	GraphiteAPIURL string
	HTTPCli        *http.Client
}

func (c *ConfigGatherer) defaults() {
	if c.HTTPCli == nil {
		c.HTTPCli = http.DefaultClient
	}
}

type gatherer struct {
	cli *graphite.Client
	cfg ConfigGatherer
}

// NewGatherer returns a new metric gatherer for graphite carbon backends.
func NewGatherer(cfg ConfigGatherer) (metric.Gatherer, error) {
	cfg.defaults()

	url, err := url.Parse(cfg.GraphiteAPIURL)
	if err != nil {
		return nil, err
	}
	cli := &graphite.Client{
		URL:    *url,
		Client: cfg.HTTPCli,
	}

	return &gatherer{
		cfg: cfg,
		cli: cli,
	}, nil
}

const (
	instantRange   = 5 * time.Minute
	targetLabelKey = "target"
)

func (g *gatherer) GatherSingle(ctx context.Context, query model.Query, t time.Time) ([]model.MetricSeries, error) {
	res, err := g.GatherRange(ctx, query, t.Add(-1*instantRange), t, 0)
	if err != nil {
		return []model.MetricSeries{}, err
	}

	if len(res) < 1 {
		return []model.MetricSeries{}, fmt.Errorf("server didn't return any metric series")
	}
	if len(res) > 1 {
		return []model.MetricSeries{}, fmt.Errorf("server returned more than one metric series, got %d", len(res))
	}

	// Get the latest datapoint.
	res[0].Metrics = res[0].Metrics[len(res[0].Metrics)-1:]

	return res, nil
}

func (g *gatherer) GatherRange(ctx context.Context, query model.Query, start, end time.Time, _ time.Duration) ([]model.MetricSeries, error) {
	// Get the data from the Graphite API.
	result, err := g.cli.QueryMulti([]string{query.Expr}, graphite.TimeInterval{
		From: start,
		To:   end,
	})
	if err != nil {
		return []model.MetricSeries{}, err
	}

	// For every metric series.
	mss := []model.MetricSeries{}
	for _, resultDps := range result {
		dps, err := resultDps.AsFloats()
		if err != nil {
			continue
		}
		if len(dps) < 1 {
			continue
		}

		// Get all it's datapoints.
		m := []model.Metric{}
		for _, dp := range dps {
			if dp.Value != nil {
				m = append(m, model.Metric{
					TS:    dp.Time,
					Value: *dp.Value,
				})
			}
		}
		ms := model.MetricSeries{
			ID: resultDps.Target,
			Labels: map[string]string{
				targetLabelKey: resultDps.Target,
			},
			Metrics: m,
		}

		mss = append(mss, ms)
	}

	return mss, nil
}
