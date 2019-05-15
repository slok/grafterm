package datasource

import (
	"context"
	"errors"
	"fmt"
	"time"

	prometheusapi "github.com/prometheus/client_golang/api"
	prometheusv1 "github.com/prometheus/client_golang/api/prometheus/v1"

	"github.com/slok/grafterm/internal/model"
	"github.com/slok/grafterm/internal/service/metric"
	"github.com/slok/grafterm/internal/service/metric/fake"
	"github.com/slok/grafterm/internal/service/metric/prometheus"
)

// ConfigGatherer is the configuration of the multi Gatherer.
type ConfigGatherer struct {
	// DashboardDatasources are the datasources that are on the dashboards and
	// will be reference, these datasources are the ones with the lowest priority.
	DashboardDatasources []model.Datasource
	// UserDatasources are the datasources outside the dashboard and defined by the suer
	// the ones that have priority over dashboards, also are the ones that will be used as
	// replacement for the aliased datasources.
	UserDatasources []model.Datasource
	// Aliases are the aliases of the dashboard datasources.
	// The key of the map is the referenced ID on the dashboard, and the
	// value of the map is the ID of the datasource that will be used.
	Aliases map[string]string
	// CreateFakeFunc is the function that will be called to create fake gatherers.
	CreateFakeFunc func(ds model.FakeDatasource) (metric.Gatherer, error)
	// CreatePrometheusFunc is the function that will be called to create Prometheus gatherers.
	CreatePrometheusFunc func(ds model.PrometheusDatasource) (metric.Gatherer, error)
}

func (c *ConfigGatherer) defaults() {
	// Set default creator function for fake.
	if c.CreateFakeFunc == nil {
		c.CreateFakeFunc = func(_ model.FakeDatasource) (metric.Gatherer, error) {
			return &fake.Gatherer{}, nil
		}
	}

	// Set default creator function for prometheus.
	if c.CreatePrometheusFunc == nil {
		c.CreatePrometheusFunc = func(ds model.PrometheusDatasource) (metric.Gatherer, error) {
			cli, err := prometheusapi.NewClient(prometheusapi.Config{
				Address: ds.Address,
			})
			if err != nil {
				return nil, err
			}
			g := prometheus.NewGatherer(prometheus.ConfigGatherer{
				Client: prometheusv1.NewAPI(cli),
			})

			return g, nil
		}
	}
	if c.Aliases == nil {
		c.Aliases = map[string]string{}
	}
}

type gatherer struct {
	cfg       ConfigGatherer
	gatherers map[string]metric.Gatherer
}

// NewGatherer returns a new gatherer that knows how to register different
// gatherer types based on datasources and when calling the methods of this
// gatherer will delegate to the correct gatherer based on the query's
// datasource ID.
func NewGatherer(cfg ConfigGatherer) (metric.Gatherer, error) {
	cfg.defaults()

	// Lowest priority (0).
	gs := map[string]metric.Gatherer{}
	for _, ds := range cfg.DashboardDatasources {
		g, err := createGatherer(cfg, ds)
		if err != nil {
			return nil, err
		}
		gs[ds.ID] = g
	}

	// Mid priority (1).
	ags := map[string]metric.Gatherer{}
	for _, ds := range cfg.UserDatasources {
		g, err := createGatherer(cfg, ds)
		if err != nil {
			return nil, err
		}
		ags[ds.ID] = g
	}

	// Use the IDs from the dashboard to use the user datasources.
	for id := range gs {
		g, ok := ags[id]
		if ok {
			gs[id] = g
		}
	}

	// Override dashboard datasource with the user datsources using the aliases.
	// Highest priority (2).
	for id, alias := range cfg.Aliases {
		ag, ok := ags[alias]
		if !ok {
			return nil, fmt.Errorf("alias %s for ID %s not found", alias, id)
		}
		gs[id] = ag
	}

	return &gatherer{
		cfg:       cfg,
		gatherers: gs,
	}, nil
}

func (g *gatherer) GatherSingle(ctx context.Context, query model.Query, t time.Time) ([]model.MetricSeries, error) {
	dsg, err := g.metricGatherer(query.DatasourceID)
	if err != nil {
		return nil, err
	}
	return dsg.GatherSingle(ctx, query, t)
}

func (g *gatherer) GatherRange(ctx context.Context, query model.Query, start, end time.Time, step time.Duration) ([]model.MetricSeries, error) {
	dsg, err := g.metricGatherer(query.DatasourceID)
	if err != nil {
		return nil, err
	}
	return dsg.GatherRange(ctx, query, start, end, step)
}

func (g *gatherer) metricGatherer(id string) (metric.Gatherer, error) {
	mg, ok := g.gatherers[id]
	if !ok {
		return nil, fmt.Errorf("datasource %s does not exists", id)
	}

	return mg, nil
}

func createGatherer(cfg ConfigGatherer, ds model.Datasource) (metric.Gatherer, error) {
	switch {
	case ds.Prometheus != nil:
		return cfg.CreatePrometheusFunc(*ds.Prometheus)
	case ds.Fake != nil:
		return cfg.CreateFakeFunc(*ds.Fake)
	}

	return nil, errors.New("not a valid datasource")
}
