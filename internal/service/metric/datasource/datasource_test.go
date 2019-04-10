package datasource_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	mmetric "github.com/slok/meterm/internal/mocks/service/metric"
	"github.com/slok/meterm/internal/model"
	"github.com/slok/meterm/internal/service/metric"
	"github.com/slok/meterm/internal/service/metric/datasource"
)

func TestGathererGatherSingle(t *testing.T) {
	datasources := []model.Datasource{
		model.Datasource{
			ID: "fakeds",
			DatasourceSource: model.DatasourceSource{
				Fake: &model.FakeDatasource{},
			},
		},
		model.Datasource{
			ID: "promds",
			DatasourceSource: model.DatasourceSource{
				Prometheus: &model.PrometheusDatasource{},
			},
		},
	}
	tests := []struct {
		name        string
		datasources []model.Datasource
		query       model.Query
		exp         func(mfake *mmetric.Gatherer, mprom *mmetric.Gatherer)
		expErr      bool
	}{
		{
			name: "A query to an non existent gatherer should fail.",
			query: model.Query{
				DatasourceID: "does-not-exists",
				Expr:         "test",
			},
			datasources: datasources,
			exp:         func(mfake *mmetric.Gatherer, mprom *mmetric.Gatherer) {},
			expErr:      true,
		},
		{
			name: "A query to the fake gatherer should use that specific gatherer.",
			query: model.Query{
				DatasourceID: "fakeds",
				Expr:         "test",
			},
			datasources: datasources,
			exp: func(mfake *mmetric.Gatherer, mprom *mmetric.Gatherer) {
				mfake.On("GatherSingle", mock.Anything, mock.Anything, mock.Anything).Once().Return([]model.MetricSeries{}, nil)
			},
		},
		{
			name: "A query to the prometheus gatherer should use that specific gatherer.",
			query: model.Query{
				DatasourceID: "promds",
				Expr:         "test",
			},
			datasources: datasources,
			exp: func(mfake *mmetric.Gatherer, mprom *mmetric.Gatherer) {
				mprom.On("GatherSingle", mock.Anything, mock.Anything, mock.Anything).Once().Return([]model.MetricSeries{}, nil)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require := require.New(t)
			assert := assert.New(t)

			// Mocks.
			mfake := &mmetric.Gatherer{}
			mprom := &mmetric.Gatherer{}
			test.exp(mfake, mprom)

			// Create the datasource based gatherer
			g, err := datasource.NewGatherer(datasource.ConfigGatherer{
				Datasources: test.datasources,

				CreateFakeFunc: func(_ model.FakeDatasource) (metric.Gatherer, error) {
					return mfake, nil
				},
				CreatePrometheusFunc: func(_ model.PrometheusDatasource) (metric.Gatherer, error) {
					return mprom, nil
				},
			})
			require.NoError(err)

			// Call to gatherer method and check the delegations to the correct
			// datasources (expected mock calls) have been made correctly.
			_, err = g.GatherSingle(context.TODO(), test.query, time.Now())
			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				mfake.AssertExpectations(t)
				mprom.AssertExpectations(t)
			}
		})
	}
}
