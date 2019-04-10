package datasource

import (
	"context"
	"errors"
	"fmt"
	"time"

	prometheusapi "github.com/prometheus/client_golang/api"
	prometheusv1 "github.com/prometheus/client_golang/api/prometheus/v1"

	"github.com/slok/meterm/internal/model"
	"github.com/slok/meterm/internal/service/metric"
	"github.com/slok/meterm/internal/service/metric/fake"
	"github.com/slok/meterm/internal/service/metric/prometheus"
)

// ConfigGatherer is the configuration of the multi Gatherer.
type ConfigGatherer struct {
	// Datasources are the configurations that the datasource gatherer
	// will use to register and create the different gatherers.
	Datasources []model.Datasource

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

	// Create the gatherers based on the datasources.
	gs := map[string]metric.Gatherer{}
	for _, ds := range cfg.Datasources {
		g, err := createGatherer(cfg, ds)
		if err != nil {
			return nil, err
		}
		gs[ds.ID] = g
	}

	return &gatherer{
		cfg:       cfg,
		gatherers: gs,
	}, nil
}

func (g *gatherer) GatherSingle(ctx context.Context, query model.Query, t time.Time) ([]model.MetricSeries, error) {
	dsg, ok := g.gatherers[query.DatasourceID]
	if !ok {
		return nil, fmt.Errorf("datasource %s does not exists", query.DatasourceID)
	}
	return dsg.GatherSingle(ctx, query, t)
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
