package fake

import (
	"context"
	"time"

	"github.com/slok/meterm/internal/model"
)

// Gatherer is a fake Gatherer.
type Gatherer struct{}

// GatherSingle satisfies metric.Gatherer interface.
func (g Gatherer) GatherSingle(_ context.Context, _ model.Query, t time.Time) ([]model.MetricSeries, error) {
	return []model.MetricSeries{
		model.MetricSeries{
			ID: "fake",
			Labels: map[string]string{
				"faked":    "true",
				"gatherer": "fake",
				"kind":     "fixed",
			},
			Metrics: []model.Metric{
				model.Metric{
					Value: float64(t.Second()),
					TS:    t,
				},
			},
		},
	}, nil
}
