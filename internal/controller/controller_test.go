package controller_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/slok/meterm/internal/controller"
	mmetric "github.com/slok/meterm/internal/mocks/service/metric"
	"github.com/slok/meterm/internal/model"
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
