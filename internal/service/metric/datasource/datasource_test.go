package datasource_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	mmetric "github.com/slok/grafterm/internal/mocks/service/metric"
	"github.com/slok/grafterm/internal/model"
	"github.com/slok/grafterm/internal/service/metric"
	"github.com/slok/grafterm/internal/service/metric/datasource"
)

func TestGathererGatherSingle(t *testing.T) {
	datasources1 := []model.Datasource{
		model.Datasource{
			ID:               "ds0",
			DatasourceSource: model.DatasourceSource{Fake: &model.FakeDatasource{}},
		},
		model.Datasource{
			ID:               "ds1",
			DatasourceSource: model.DatasourceSource{Prometheus: &model.PrometheusDatasource{}},
		},
	}
	datasources2 := []model.Datasource{
		model.Datasource{
			ID:               "ds2",
			DatasourceSource: model.DatasourceSource{Graphite: &model.GraphiteDatasource{}},
		},
		model.Datasource{
			ID:               "ds3",
			DatasourceSource: model.DatasourceSource{Prometheus: &model.PrometheusDatasource{}},
		},
	}
	datasources3 := []model.Datasource{
		model.Datasource{
			ID:               "ds0",
			DatasourceSource: model.DatasourceSource{Fake: &model.FakeDatasource{}},
		},
		model.Datasource{
			ID:               "ds1",
			DatasourceSource: model.DatasourceSource{Prometheus: &model.PrometheusDatasource{}},
		},
	}
	tests := []struct {
		name                 string
		dashboardDatasources []model.Datasource
		userDatasources      []model.Datasource
		aliases              map[string]string
		query                model.Query
		exp                  func(mgs []*mmetric.Gatherer)
		expErr               bool
	}{
		{
			name: "A query to an non existent gatherer should fail.",
			query: model.Query{
				DatasourceID: "does-not-exists",
				Expr:         "test",
			},
			dashboardDatasources: datasources1,
			userDatasources:      datasources2,
			exp:                  func(mgs []*mmetric.Gatherer) {},
			expErr:               true,
		},
		{
			name: "A query using a datasource ID should use the correct datasource based on the query.",
			query: model.Query{
				DatasourceID: "ds1",
				Expr:         "test",
			},
			dashboardDatasources: datasources1,
			userDatasources:      datasources2,
			exp: func(mgs []*mmetric.Gatherer) {
				mgs[1].On("GatherSingle", mock.Anything, mock.Anything, mock.Anything).Once().Return([]model.MetricSeries{}, nil)
			},
		},
		{
			name: "A query using a datasource ID that isn't on the dashboard datasources but is on the user datasources should fail.",
			query: model.Query{
				DatasourceID: "ds3",
				Expr:         "test",
			},
			dashboardDatasources: datasources1,
			userDatasources:      datasources2,
			exp:                  func(mgs []*mmetric.Gatherer) {},
			expErr:               true,
		},
		{
			name: "A query using a datasource ID that is aliased should use the alias replacement datasource.",
			query: model.Query{
				DatasourceID: "ds1",
				Expr:         "test",
			},
			dashboardDatasources: datasources1,
			userDatasources:      datasources2,
			aliases: map[string]string{
				"ds1": "ds3",
			},
			exp: func(mgs []*mmetric.Gatherer) {
				mgs[3].On("GatherSingle", mock.Anything, mock.Anything, mock.Anything).Once().Return([]model.MetricSeries{}, nil)
			},
		},
		{
			name: "A query using a datasource ID that has the same ID on the user defined datasources, should use the user one.",
			query: model.Query{
				DatasourceID: "ds1",
				Expr:         "test",
			},
			dashboardDatasources: datasources1,
			userDatasources:      datasources3,
			exp: func(mgs []*mmetric.Gatherer) {
				mgs[3].On("GatherSingle", mock.Anything, mock.Anything, mock.Anything).Once().Return([]model.MetricSeries{}, nil)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require := require.New(t)
			assert := assert.New(t)

			// Create mocks based on the datasources of the test.
			mgs := []*mmetric.Gatherer{}
			for i := 0; i < len(test.dashboardDatasources); i++ {
				mgs = append(mgs, &mmetric.Gatherer{})
			}
			for i := 0; i < len(test.userDatasources); i++ {
				mgs = append(mgs, &mmetric.Gatherer{})
			}
			test.exp(mgs)

			// Create the datasource based gatherer.
			// The creation funcs return the mocks in order.
			gCount := 0
			g, err := datasource.NewGatherer(datasource.ConfigGatherer{
				DashboardDatasources: test.dashboardDatasources,
				UserDatasources:      test.userDatasources,
				Aliases:              test.aliases,
				CreateFakeFunc: func(_ model.FakeDatasource) (metric.Gatherer, error) {
					g := mgs[gCount]
					gCount++
					return g, nil
				},
				CreatePrometheusFunc: func(_ model.PrometheusDatasource) (metric.Gatherer, error) {
					g := mgs[gCount]
					gCount++
					return g, nil
				},
				CreateGraphiteFunc: func(_ model.GraphiteDatasource) (metric.Gatherer, error) {
					g := mgs[gCount]
					gCount++
					return g, nil
				},
			})
			require.NoError(err)

			// Call to gatherer method and check the delegations to the correct
			// datasources (expected mock calls) have been made correctly.
			_, err = g.GatherSingle(context.TODO(), test.query, time.Now())
			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				for _, mg := range mgs {
					mg.AssertExpectations(t)
				}
			}
		})
	}
}
