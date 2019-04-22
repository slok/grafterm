package controller_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/slok/grafterm/internal/controller"
	mmetric "github.com/slok/grafterm/internal/mocks/service/metric"
	"github.com/slok/grafterm/internal/model"
)

func TestGetSingleMetric(t *testing.T) {
	tests := []struct {
		name           string
		query          model.Query
		serviceMetrics []model.MetricSeries
		serviceErr     error
		ts             time.Time
		expErr         bool
		expMetric      *model.Metric
	}{
		{
			name:  "Returning a correct metric the controller should handle the single metric correctly.",
			query: model.Query{Expr: "test"},
			serviceMetrics: []model.MetricSeries{
				model.MetricSeries{
					Metrics: []model.Metric{
						{Value: 17.9},
					},
				},
			},
			ts:        time.Now(),
			expMetric: &model.Metric{Value: 17.9},
		},
		{
			name:  "Returning multiple metrics should error.",
			query: model.Query{Expr: "test"},
			serviceMetrics: []model.MetricSeries{
				model.MetricSeries{
					Metrics: []model.Metric{
						{Value: 17.9},
						{Value: 28.1},
					},
				},
			},
			ts:     time.Now(),
			expErr: true,
		},
		{
			name:  "Returning no metrics should error.",
			query: model.Query{Expr: "test"},
			serviceMetrics: []model.MetricSeries{
				model.MetricSeries{
					Metrics: []model.Metric{},
				},
			},
			ts:     time.Now(),
			expErr: true,
		},
		{
			name:           "Returning no metric series should error.",
			query:          model.Query{Expr: "test"},
			serviceMetrics: []model.MetricSeries{},
			ts:             time.Now(),
			expErr:         true,
		},
		{
			name:  "Returning multiple metric series should error.",
			query: model.Query{Expr: "test"},
			serviceMetrics: []model.MetricSeries{
				model.MetricSeries{},
				model.MetricSeries{},
			},
			ts:     time.Now(),
			expErr: true,
		},
		{
			name:       "Returning a error from the metrics service should error.",
			query:      model.Query{Expr: "test"},
			serviceErr: errors.New("wanted error"),
			ts:         time.Now(),
			expErr:     true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			// Mocks.
			mg := &mmetric.Gatherer{}
			mg.On("GatherSingle", mock.Anything, test.query, test.ts).Once().Return(test.serviceMetrics, test.serviceErr)

			c := controller.NewController(mg)
			gotm, err := c.GetSingleMetric(context.TODO(), test.query, test.ts)

			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				assert.Equal(test.expMetric, gotm)
				mg.AssertExpectations(t)
			}
		})
	}
}

func TestGetRangeMetrics(t *testing.T) {
	start := time.Now()
	end := start.Add(5 * time.Hour)

	tests := []struct {
		name           string
		query          model.Query
		serviceMetrics []model.MetricSeries
		serviceErr     error
		start          time.Time
		end            time.Time
		step           time.Duration
		expErr         bool
		expSeries      []model.MetricSeries
	}{
		{
			name:   "Receiving a non positive step should return an error.",
			query:  model.Query{Expr: "test"},
			start:  start,
			end:    end,
			expErr: true,
		},
		{
			name:   "Receiving a start TS that is older than a end TS should return an error.",
			query:  model.Query{Expr: "test"},
			start:  end,
			end:    start,
			step:   1 * time.Hour,
			expErr: true,
		},
		{
			name:  "Receiving and error from the services should return an error.",
			query: model.Query{Expr: "test"},
			serviceMetrics: []model.MetricSeries{
				model.MetricSeries{
					Metrics: []model.Metric{
						{Value: 17.9},
					},
				},
			},
			serviceErr: errors.New("wanted error"),
			start:      start,
			end:        end,
			step:       1 * time.Hour,
			expErr:     true,
		},
		{
			name:  "Receiving correct group of arguments should call the services and return the metrics.",
			query: model.Query{Expr: "test"},
			serviceMetrics: []model.MetricSeries{
				model.MetricSeries{
					Metrics: []model.Metric{
						{Value: 17.9},
					},
				},
			},
			start: start,
			end:   end,
			step:  1 * time.Hour,
			expSeries: []model.MetricSeries{
				model.MetricSeries{
					Metrics: []model.Metric{
						{Value: 17.9},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			// Mocks.
			mg := &mmetric.Gatherer{}
			mg.On("GatherRange", mock.Anything, test.query, test.start, test.end, test.step).Once().Return(test.serviceMetrics, test.serviceErr)

			c := controller.NewController(mg)
			gotSeries, err := c.GetRangeMetrics(context.TODO(), test.query, test.start, test.end, test.step)

			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				assert.Equal(test.expSeries, gotSeries)
				mg.AssertExpectations(t)
			}
		})
	}
}
