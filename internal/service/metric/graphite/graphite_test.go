package graphite_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/slok/grafterm/internal/model"
	"github.com/slok/grafterm/internal/service/metric/graphite"
)

func TestGathererGatherSingle(t *testing.T) {
	tests := map[string]struct {
		graphiteResponse string
		cfg              graphite.ConfigGatherer
		expMetricSeries  []model.MetricSeries
		expErr           bool
	}{
		"Getting multiple metrics in a single metric series should return only the last metric.": {
			graphiteResponse: `
[{
  "target": "batman",
  "datapoints": [
	[612.54, 1558275625],
	[712.54, 1558275725],
	[812.54, 1558275825],
	[912.54, 1558275925]
  ]
}]`,
			expMetricSeries: []model.MetricSeries{
				{
					ID:     "batman",
					Labels: map[string]string{"target": "batman"},
					Metrics: []model.Metric{
						{
							Value: 912.54,
							TS:    time.Unix(1558275925, 0),
						},
					},
				},
			},
		},
		"Getting 0 metric series should error.": {
			graphiteResponse: `[]`,
			expErr:           true,
		},
		"Getting more than one metric series should error.": {
			graphiteResponse: `
[
	{"target": "batman","datapoints": [[612.54, 1558275625]]},
	{"target2": "batman","datapoints": [[612.54, 1558275625]]},
	{"target3": "batman","datapoints": [[612.54, 1558275625]]}
]`,
			expErr: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			// Mock server response.
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(test.graphiteResponse))
			}))
			defer srv.Close()
			test.cfg.GraphiteAPIURL = srv.URL

			g, _ := graphite.NewGatherer(test.cfg)
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
		graphiteResponse string
		cfg              graphite.ConfigGatherer
		expMetricSeries  []model.MetricSeries
		expErr           bool
	}{
		"When Graphite API returns multiple time series the gatherer should return the metrics translated to the model.": {
			graphiteResponse: `
[
	{
		"target": "batman",
		"datapoints": [[612.54, 1558332722],[712.54, 1558332724],[812.54, 1558332726],[912.54, 1558332728]]
	},
	{
		"target": "deadpool",
		"datapoints": [[10.15, 1558332822],[10.17, 1558332832],[13.7, 1558332843]]
	},
	{
		"target": "wolverine",
		"datapoints": [[10012.8992, 1558352722],[60072.8992, 1558352726]]
	}
]`,
			expMetricSeries: []model.MetricSeries{
				{
					ID:     "batman",
					Labels: map[string]string{"target": "batman"},
					Metrics: []model.Metric{
						{Value: 612.54, TS: time.Unix(1558332722, 0)},
						{Value: 712.54, TS: time.Unix(1558332724, 0)},
						{Value: 812.54, TS: time.Unix(1558332726, 0)},
						{Value: 912.54, TS: time.Unix(1558332728, 0)},
					},
				},
				{
					ID:     "deadpool",
					Labels: map[string]string{"target": "deadpool"},
					Metrics: []model.Metric{
						{Value: 10.15, TS: time.Unix(1558332822, 0)},
						{Value: 10.17, TS: time.Unix(1558332832, 0)},
						{Value: 13.7, TS: time.Unix(1558332843, 0)},
					},
				},
				{
					ID:     "wolverine",
					Labels: map[string]string{"target": "wolverine"},
					Metrics: []model.Metric{
						{Value: 10012.8992, TS: time.Unix(1558352722, 0)},
						{Value: 60072.8992, TS: time.Unix(1558352726, 0)},
					},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			// Mock server response.
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(test.graphiteResponse))
			}))
			defer srv.Close()
			test.cfg.GraphiteAPIURL = srv.URL

			g, _ := graphite.NewGatherer(test.cfg)
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
