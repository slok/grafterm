package prometheus_test

import (
	"context"
	"errors"
	"testing"
	"time"

	prommodel "github.com/prometheus/common/model"
	mpromv1 "github.com/slok/grafterm/internal/mocks/github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/slok/grafterm/internal/model"
	"github.com/slok/grafterm/internal/service/metric/prometheus"
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
						"k2":       "v2",
						"__name__": "test-metric",
					},
					Value:     prommodel.SampleValue(1.2),
					Timestamp: prommodel.TimeFromUnixNano(now.UnixNano()),
				},
			},
			expMetricSeries: []model.MetricSeries{
				model.MetricSeries{
					ID: `test-metric{k1="v1", k2="v2"}`,
					Labels: map[string]string{
						"k1":       "v1",
						"k2":       "v2",
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
			mapi.On("Query", mock.Anything, mock.Anything, mock.Anything).Once().Return(test.prommetric, nil, expErr)
			test.cfg.Client = mapi

			g := prometheus.NewGatherer(test.cfg)
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
	t1 := time.Now().Truncate(time.Second).Add(-1 * time.Minute)
	t2 := time.Now().Truncate(time.Second).Add(-2 * time.Minute)
	t3 := time.Now().Truncate(time.Second).Add(-3 * time.Minute)
	t4 := time.Now().Truncate(time.Second).Add(-4 * time.Minute)
	t5 := time.Now().Truncate(time.Second).Add(-5 * time.Minute)
	t6 := time.Now().Truncate(time.Second).Add(-6 * time.Minute)
	t7 := time.Now().Truncate(time.Second).Add(-7 * time.Minute)

	tests := []struct {
		name            string
		prommetric      prommodel.Value
		cfg             prometheus.ConfigGatherer
		expMetricSeries []model.MetricSeries
		expErr          bool
	}{
		{
			name: "When Prometheus returns a matrix metrics the Gatherer should return the translated metric.",
			prommetric: prommodel.Matrix{
				&prommodel.SampleStream{
					Metric: prommodel.Metric{"k1": "v1", "k2": "v2", "__name__": "test-metric"},
					Values: []prommodel.SamplePair{
						{
							Value:     prommodel.SampleValue(1.2),
							Timestamp: prommodel.TimeFromUnixNano(t1.UnixNano()),
						},
						{
							Value:     prommodel.SampleValue(2.3),
							Timestamp: prommodel.TimeFromUnixNano(t2.UnixNano()),
						},
						{
							Value:     prommodel.SampleValue(3.4),
							Timestamp: prommodel.TimeFromUnixNano(t3.UnixNano()),
						},
						{
							Value:     prommodel.SampleValue(4.5),
							Timestamp: prommodel.TimeFromUnixNano(t4.UnixNano()),
						},
						{
							Value:     prommodel.SampleValue(5.6),
							Timestamp: prommodel.TimeFromUnixNano(t5.UnixNano()),
						},
					},
				},
				&prommodel.SampleStream{
					Metric: prommodel.Metric{"k1": "v1", "k2": "v2", "k3": "v3", "__name__": "test-metric"},
					Values: []prommodel.SamplePair{
						{
							Value:     prommodel.SampleValue(6.7),
							Timestamp: prommodel.TimeFromUnixNano(t1.UnixNano()),
						},
						{
							Value:     prommodel.SampleValue(7.8),
							Timestamp: prommodel.TimeFromUnixNano(t2.UnixNano()),
						},
						{
							Value:     prommodel.SampleValue(9.10),
							Timestamp: prommodel.TimeFromUnixNano(t3.UnixNano()),
						},
					},
				},
				&prommodel.SampleStream{
					Metric: prommodel.Metric{"k5": "v5", "__name__": "test-metric"},
					Values: []prommodel.SamplePair{
						{
							Value:     prommodel.SampleValue(10.11),
							Timestamp: prommodel.TimeFromUnixNano(t6.UnixNano()),
						},
					},
				},
				&prommodel.SampleStream{
					Metric: prommodel.Metric{"k1": "v1", "k2": "v2", "__name__": "test-metric2"},
					Values: []prommodel.SamplePair{
						{
							Value:     prommodel.SampleValue(11.12),
							Timestamp: prommodel.TimeFromUnixNano(t7.UnixNano()),
						},
					},
				},
			},
			expMetricSeries: []model.MetricSeries{
				model.MetricSeries{
					ID:     `test-metric{k1="v1", k2="v2"}`,
					Labels: map[string]string{"k1": "v1", "k2": "v2", "__name__": "test-metric"},
					Metrics: []model.Metric{
						model.Metric{Value: 1.2, TS: t1},
						model.Metric{Value: 2.3, TS: t2},
						model.Metric{Value: 3.4, TS: t3},
						model.Metric{Value: 4.5, TS: t4},
						model.Metric{Value: 5.6, TS: t5},
					},
				},
				model.MetricSeries{
					ID:     `test-metric{k1="v1", k2="v2", k3="v3"}`,
					Labels: map[string]string{"k1": "v1", "k2": "v2", "k3": "v3", "__name__": "test-metric"},
					Metrics: []model.Metric{
						model.Metric{Value: 6.7, TS: t1},
						model.Metric{Value: 7.8, TS: t2},
						model.Metric{Value: 9.10, TS: t3},
					},
				},
				model.MetricSeries{
					ID:     `test-metric{k5="v5"}`,
					Labels: map[string]string{"k5": "v5", "__name__": "test-metric"},
					Metrics: []model.Metric{
						model.Metric{Value: 10.11, TS: t6},
					},
				},
				model.MetricSeries{
					ID:     `test-metric2{k1="v1", k2="v2"}`,
					Labels: map[string]string{"k1": "v1", "k2": "v2", "__name__": "test-metric2"},
					Metrics: []model.Metric{
						model.Metric{Value: 11.12, TS: t7},
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
			mapi.On("QueryRange", mock.Anything, mock.Anything, mock.Anything).Once().Return(test.prommetric, nil, expErr)
			test.cfg.Client = mapi

			g := prometheus.NewGatherer(test.cfg)
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
