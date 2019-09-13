package influxdb_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	_ "github.com/influxdata/influxdb1-client" // needed due to go mod bug
	influxdbv2 "github.com/influxdata/influxdb1-client/v2"

	"github.com/stretchr/testify/assert"

	"github.com/slok/grafterm/internal/model"
	"github.com/slok/grafterm/internal/service/metric/influxdb"
)

func TestGathererGatherSingle(t *testing.T) {
	tests := map[string]struct {
		influxdbResponse string
		cfg              influxdb.ConfigGatherer
		expMetricSeries  []model.MetricSeries
		expErr           bool
	}{
		"Getting multiple metrics in a single metric series should return only the last metric": {
			influxdbResponse: `
{"results":[
  {
  "series": [
     {"name":"myseries","columns":["time","mean"],"values":[["2017-03-01T00:16:18Z",12.34],["2017-03-01T00:17:18Z",23.45],["2017-03-01T00:18:18Z",34.56],["2017-03-01T00:19:18Z",45.67]]}
  ]
  }
]}`,
			expMetricSeries: []model.MetricSeries{
				{
					ID: "myseries",
					Metrics: []model.Metric{
						{
							Value: 45.67,
							TS:    time.Date(2017, 3, 1, 0, 19, 18, 0, time.UTC),
						},
					},
				},
			},
		},
		"Getting 0 metric series should error": {
			influxdbResponse: `[]`,
			expErr:           true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			// Mock server response.
			srv := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.Write([]byte(test.influxdbResponse))
				}))
			defer srv.Close()

			test.cfg.Addr = srv.URL
			test.cfg.Client = influxdbClient(srv.URL)
			test.cfg.Database = "dummy"

			g, _ := influxdb.NewGatherer(test.cfg)
			gotms, err := g.GatherSingle(context.TODO(), model.Query{}, time.Now())
			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				// We don't control the order of the MetricSeries and sorting is harder than checking in
				// two steps.
				assert.Len(gotms, len(test.expMetricSeries))
				for _, gotm := range gotms {
					assert.Contains(test.expMetricSeries, gotm)
				}
			}
		})
	}
}

func TestGathererGatherRange(t *testing.T) {
	tests := map[string]struct {
		influxdbResponse string
		cfg              influxdb.ConfigGatherer
		expMetricSeries  []model.MetricSeries
		expErr           bool
	}{
		"When influxdb returns multiple time series the gatherer should return the metrics translated to the model": {
			influxdbResponse: `
{"results":[
  {
  "series": [
     {"name":"myseries1","columns":["time","mean1"],"values":[["2017-03-01T00:16:18Z",12.34],["2017-03-01T00:17:18Z",23.45]]},
     {"name":"myseries2","columns":["time","mean2"],"values":[["2017-03-01T00:18:18Z",34.56],["2017-03-01T00:19:18Z",45.67]]}
  ]
  }
]}`,
			expMetricSeries: []model.MetricSeries{
				{
					ID: "myseries1",
					Metrics: []model.Metric{
						{Value: 12.34, TS: time.Date(2017, 3, 1, 0, 16, 18, 0, time.UTC)},
						{Value: 23.45, TS: time.Date(2017, 3, 1, 0, 17, 18, 0, time.UTC)},
					},
				},
				{
					ID: "myseries2",
					Metrics: []model.Metric{
						{Value: 34.56, TS: time.Date(2017, 3, 1, 0, 18, 18, 0, time.UTC)},
						{Value: 45.67, TS: time.Date(2017, 3, 1, 0, 19, 18, 0, time.UTC)},
					},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			// Mock server response.
			srv := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.Write([]byte(test.influxdbResponse))
				}))
			defer srv.Close()

			test.cfg.Addr = srv.URL
			test.cfg.Client = influxdbClient(srv.URL)
			test.cfg.Database = "dummy"

			g, _ := influxdb.NewGatherer(test.cfg)
			gotms, err := g.GatherRange(context.TODO(), model.Query{}, time.Now(), time.Now(), 0)
			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				// We don't control the order of the MetricSeries and sorting is harder than checking in
				// two steps.
				assert.Len(gotms, len(test.expMetricSeries))
				for _, gotm := range gotms {
					assert.Contains(test.expMetricSeries, gotm)
				}
			}
		})
	}
}

func influxdbClient(addr string) influxdbv2.Client {
	cli, _ := influxdbv2.NewHTTPClient(
		influxdbv2.HTTPConfig{
			Addr: addr,
		},
	)
	return cli
}
