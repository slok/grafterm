package prometheus

import (
	"context"
	"errors"
	"strings"
	"time"

	promv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	prommodel "github.com/prometheus/common/model"

	"github.com/slok/grafterm/internal/model"
	"github.com/slok/grafterm/internal/service/metric"
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
	val, _, err := g.cli.Query(ctx, query.Expr, t)
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

func (g *gatherer) GatherRange(ctx context.Context, query model.Query, start, end time.Time, step time.Duration) ([]model.MetricSeries, error) {
	// Get value from Prometheus.
	val, _, err := g.cli.QueryRange(ctx, query.Expr, promv1.Range{
		Start: start,
		End:   end,
		Step:  step,
	})
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
		matrix := pm.(prommodel.Matrix)
		res = g.transformMatrix(matrix)
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
// A Prometheus vector is an slice of metrics (group of labels) that have one
// sample only (all samples from all metrics have the same timestamp)
func (g *gatherer) transformVector(vector prommodel.Vector) []model.MetricSeries {
	res := []model.MetricSeries{}

	for _, sample := range vector {
		id := sample.Metric.String()
		labels := g.labelSetToMap(prommodel.LabelSet(sample.Metric))
		series := model.MetricSeries{
			ID:     id,
			Labels: g.sanitizeLabels(labels),
		}

		// Add the metric to the series.
		series.Metrics = append(series.Metrics, model.Metric{
			TS:    sample.Timestamp.Time(),
			Value: float64(sample.Value),
		})

		res = append(res, series)
	}

	return res
}

// transformMatrix will get a prometheus Matrix and transform to a domain model
// MetricSeries slice.
// A Prometheus Matrix is an slices of metrics (group of labels) that have multiple
// samples (in a slice of samples).
func (g *gatherer) transformMatrix(matrix prommodel.Matrix) []model.MetricSeries {
	res := []model.MetricSeries{}

	// Use a map to index the different series based on labels.
	for _, sampleStream := range matrix {
		id := sampleStream.Metric.String()
		labels := g.labelSetToMap(prommodel.LabelSet(sampleStream.Metric))
		series := model.MetricSeries{
			ID:     id,
			Labels: g.sanitizeLabels(labels),
		}

		// Add the metric to the series.
		for _, sample := range sampleStream.Values {
			series.Metrics = append(series.Metrics, model.Metric{
				TS:    sample.Timestamp.Time(),
				Value: float64(sample.Value),
			})
		}

		res = append(res, series)
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
