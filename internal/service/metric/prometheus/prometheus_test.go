package prometheus_test

import (
	"context"
	"errors"
	"testing"
	"time"

	prommodel "github.com/prometheus/common/model"
	mpromv1 "github.com/slok/meterm/internal/mocks/github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/slok/meterm/internal/model"
	"github.com/slok/meterm/internal/service/log"
	"github.com/slok/meterm/internal/service/metric/prometheus"
)

func TestGathererGatherSingle(t *testing.T) {
	now := time.Now().Truncate(time.Second)

	tests := []struct {
		name            string
		prommetric      prommodel.Value
		cfg             prometheus.ConfigGatherer
		expMetricSeries []model.MetricSeries
		expErr          bool
	}{
		{
			name: "When Prometheus returns a scalar the Gatherer should return a metric.",
			prommetric: &prommodel.Scalar{
				Timestamp: prommodel.TimeFromUnixNano(now.UnixNano()),
				Value:     prommodel.SampleValue(4.6),
			},
			expMetricSeries: []model.MetricSeries{
				model.MetricSeries{
					Metrics: []model.Metric{
						model.Metric{
							Value: 4.6,
							TS:    now,
						},
					},
				},
			},
		},
		{
			name: "When Prometheus returns a Vector with single metric the Gatherer should return the translated metric.",
			prommetric: prommodel.Vector{
				&prommodel.Sample{
					Metric: prommodel.Metric{
						"k1":       "v1",
						"k2":       "k2",
						"__name__": "test-metric",
					},
					Value:     prommodel.SampleValue(1.2),
					Timestamp: prommodel.TimeFromUnixNano(now.UnixNano()),
				},
			},
			expMetricSeries: []model.MetricSeries{
				model.MetricSeries{
					ID: "test-metric",
					Labels: map[string]string{
						"k1":       "v1",
						"k2":       "k2",
						"__name__": "test-metric",
					},
					Metrics: []model.Metric{
						model.Metric{
							Value: 1.2,
							TS:    now,
						},
					},
				},
			},
		},
		{
			name: "When Prometheus returns a Vector with multiple metrics the Gatherer should return the translated metric.",
			prommetric: prommodel.Vector{
				&prommodel.Sample{
					Metric:    prommodel.Metric{"k1": "v1", "k2": "k2", "__name__": "test-metric"},
					Value:     prommodel.SampleValue(1.2),
					Timestamp: prommodel.TimeFromUnixNano(now.UnixNano()),
				},
				&prommodel.Sample{
					Metric:    prommodel.Metric{"k1": "v1", "k2": "k2", "__name__": "test-metric"},
					Value:     prommodel.SampleValue(3.4),
					Timestamp: prommodel.TimeFromUnixNano(now.UnixNano()),
				},
				&prommodel.Sample{
					Metric:    prommodel.Metric{"k1": "v1", "k2": "k2", "__name__": "test-metric"},
					Value:     prommodel.SampleValue(5.5),
					Timestamp: prommodel.TimeFromUnixNano(now.UnixNano()),
				},
				&prommodel.Sample{
					Metric:    prommodel.Metric{"k1": "v1", "k2": "k2", "k3": "k3", "__name__": "test-metric"},
					Value:     prommodel.SampleValue(6.7),
					Timestamp: prommodel.TimeFromUnixNano(now.UnixNano()),
				},
				&prommodel.Sample{
					Metric:    prommodel.Metric{"k5": "v5", "__name__": "test-metric"},
					Value:     prommodel.SampleValue(8.9),
					Timestamp: prommodel.TimeFromUnixNano(now.UnixNano()),
				},
				&prommodel.Sample{
					Metric:    prommodel.Metric{"k1": "v1", "k2": "k2", "__name__": "test-metric2"},
					Value:     prommodel.SampleValue(10.11),
					Timestamp: prommodel.TimeFromUnixNano(now.UnixNano()),
				},
			},
			expMetricSeries: []model.MetricSeries{
				model.MetricSeries{
					ID:     "test-metric",
					Labels: map[string]string{"k1": "v1", "k2": "k2", "__name__": "test-metric"},
					Metrics: []model.Metric{
						model.Metric{Value: 1.2, TS: now},
						model.Metric{Value: 3.4, TS: now},
						model.Metric{Value: 5.5, TS: now},
					},
				},
				model.MetricSeries{
					ID:     "test-metric",
					Labels: map[string]string{"k1": "v1", "k2": "k2", "k3": "k3", "__name__": "test-metric"},
					Metrics: []model.Metric{
						model.Metric{Value: 6.7, TS: now},
					},
				},
				model.MetricSeries{
					ID:     "test-metric",
					Labels: map[string]string{"k5": "v5", "__name__": "test-metric"},
					Metrics: []model.Metric{
						model.Metric{Value: 8.9, TS: now},
					},
				},
				model.MetricSeries{
					ID:     "test-metric2",
					Labels: map[string]string{"k1": "v1", "k2": "k2", "__name__": "test-metric2"},
					Metrics: []model.Metric{
						model.Metric{Value: 10.11, TS: now},
					},
				},
			},
		},
		{
			name: "When Prometheus returns a Vector with multiple metrics with special filter the Gatherer should return the translated metric.",
			cfg: prometheus.ConfigGatherer{
				FilterSpecialLabels: true,
			},
			prommetric: prommodel.Vector{
				&prommodel.Sample{
					Metric:    prommodel.Metric{"k1": "v1", "k2": "k2", "__name__": "test-metric"},
					Value:     prommodel.SampleValue(1.2),
					Timestamp: prommodel.TimeFromUnixNano(now.UnixNano()),
				},
				&prommodel.Sample{
					Metric:    prommodel.Metric{"k1": "v1", "k2": "k2", "__name__": "test-metric"},
					Value:     prommodel.SampleValue(3.4),
					Timestamp: prommodel.TimeFromUnixNano(now.UnixNano()),
				},
				&prommodel.Sample{
					Metric:    prommodel.Metric{"k1": "v1", "k2": "k2", "__name__": "test-metric"},
					Value:     prommodel.SampleValue(5.5),
					Timestamp: prommodel.TimeFromUnixNano(now.UnixNano()),
				},
				&prommodel.Sample{
					Metric:    prommodel.Metric{"k1": "v1", "k2": "k2", "k3": "k3", "__name__": "test-metric"},
					Value:     prommodel.SampleValue(6.7),
					Timestamp: prommodel.TimeFromUnixNano(now.UnixNano()),
				},
				&prommodel.Sample{
					Metric:    prommodel.Metric{"k5": "v5", "__name__": "test-metric"},
					Value:     prommodel.SampleValue(8.9),
					Timestamp: prommodel.TimeFromUnixNano(now.UnixNano()),
				},
				&prommodel.Sample{
					Metric:    prommodel.Metric{"k1": "v1", "k2": "k2", "__name__": "test-metric2"},
					Value:     prommodel.SampleValue(10.11),
					Timestamp: prommodel.TimeFromUnixNano(now.UnixNano()),
				},
			},
			expMetricSeries: []model.MetricSeries{
				model.MetricSeries{
					ID:     "test-metric",
					Labels: map[string]string{"k1": "v1", "k2": "k2"},
					Metrics: []model.Metric{
						model.Metric{Value: 1.2, TS: now},
						model.Metric{Value: 3.4, TS: now},
						model.Metric{Value: 5.5, TS: now},
					},
				},
				model.MetricSeries{
					ID:     "test-metric",
					Labels: map[string]string{"k1": "v1", "k2": "k2", "k3": "k3"},
					Metrics: []model.Metric{
						model.Metric{Value: 6.7, TS: now},
					},
				},
				model.MetricSeries{
					ID:     "test-metric",
					Labels: map[string]string{"k5": "v5"},
					Metrics: []model.Metric{
						model.Metric{Value: 8.9, TS: now},
					},
				},
				model.MetricSeries{
					ID:     "test-metric2",
					Labels: map[string]string{"k1": "v1", "k2": "k2"},
					Metrics: []model.Metric{
						model.Metric{Value: 10.11, TS: now},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			var expErr error
			if test.expErr {
				expErr = errors.New("wanted error")
			}

			// Mocks.
			mapi := &mpromv1.API{}
			mapi.On("Query", mock.Anything, mock.Anything, mock.Anything).Once().Return(test.prommetric, expErr)
			test.cfg.Client = mapi

			g := prometheus.NewGatherer(test.cfg, log.Dummy)
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
