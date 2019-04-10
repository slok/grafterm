package prometheus

import (
	"context"
	"errors"
	"strings"
	"time"

	promv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	prommodel "github.com/prometheus/common/model"

	"github.com/slok/meterm/internal/model"
	"github.com/slok/meterm/internal/service/metric"
)

// ConfigGatherer is the configuration of the Prometheus gatherer.
type ConfigGatherer struct {
	// Client is the prometheus API client.
	Client promv1.API
	// FilterSpecialLabels will return the metrics with the special labels filtered.
	// The special labels start with `__`, examples: `__name__`, `__scheme__`.
	FilterSpecialLabels bool
}

type gatherer struct {
	cli promv1.API
	cfg ConfigGatherer
}

// NewGatherer returns a new metric gatherer for prometheus backends.
func NewGatherer(cfg ConfigGatherer) metric.Gatherer {
	return &gatherer{
		cli: cfg.Client,
		cfg: cfg,
	}
}

func (g *gatherer) GatherSingle(ctx context.Context, query model.Query, t time.Time) ([]model.MetricSeries, error) {
	// Get value from Prometheus.
	val, err := g.cli.Query(ctx, query.Expr, t)
	if err != nil {
		return []model.MetricSeries{}, err
	}

	// Translate prom values to domain.
	res, err := g.promToModel(val)
	if err != nil {
		return []model.MetricSeries{}, err
	}

	return res, nil
}

// promToModel converts a prometheus result metric to a domain model one.
func (g *gatherer) promToModel(pm prommodel.Value) ([]model.MetricSeries, error) {
	res := []model.MetricSeries{}

	switch pm.Type() {
	case prommodel.ValScalar:
		scalar := pm.(*prommodel.Scalar)
		res = g.transformScalar(scalar)
	case prommodel.ValVector:
		vector := pm.(prommodel.Vector)
		res = g.transformVector(vector)
	case prommodel.ValMatrix:
		// TODO(slok).
	default:
		return res, errors.New("prometheus value type not supported")
	}

	return res, nil
}

// transformScalar will get a prometheus Scalar and transform to a domain model
// MetricSeries slice.
func (g *gatherer) transformScalar(scalar *prommodel.Scalar) []model.MetricSeries {
	res := []model.MetricSeries{}

	m := model.Metric{
		TS:    scalar.Timestamp.Time(),
		Value: float64(scalar.Value),
	}
	res = append(res, model.MetricSeries{
		Metrics: []model.Metric{m},
	})

	return res
}

// transformVector will get a prometheus Vector and transform to a domain model
// MetricSeries slice.
func (g *gatherer) transformVector(vector prommodel.Vector) []model.MetricSeries {
	res := []model.MetricSeries{}

	// Use a map to index the different series based on labels.
	indexedSeries := map[string]model.MetricSeries{}
	for _, sample := range vector {
		index := sample.Metric.String()

		// Do we already have the series? If not create a new one.
		series, ok := indexedSeries[index]
		if !ok {
			labels := g.labelSetToMap(prommodel.LabelSet(sample.Metric))
			series = model.MetricSeries{
				ID:     g.getMetricName(labels),
				Labels: g.sanitizeLabels(labels),
			}
			indexedSeries[index] = series
		}

		// Add the metric to the series.
		series.Metrics = append(series.Metrics, model.Metric{
			TS:    sample.Timestamp.Time(),
			Value: float64(sample.Value),
		})
		indexedSeries[index] = series
	}

	for _, v := range indexedSeries {
		res = append(res, v)
	}

	return res
}

func (g *gatherer) labelSetToMap(ls prommodel.LabelSet) map[string]string {
	res := map[string]string{}
	for k, v := range ls {
		res[string(k)] = string(v)
	}

	return res
}

// sanitizeLabels will sanitize the map label values.
// 	- Remove special labels if required (start with `__`).
func (g *gatherer) sanitizeLabels(m map[string]string) map[string]string {

	// Filter if required.
	if !g.cfg.FilterSpecialLabels {
		return m
	}

	res := map[string]string{}
	for k, v := range m {
		if strings.HasPrefix(k, "__") {
			continue
		}
		res[k] = v
	}

	return res
}

func (g *gatherer) getMetricName(labels map[string]string) string {
	id, ok := labels[prommodel.MetricNameLabel]
	if !ok {
		return ""
	}
	return id
}
