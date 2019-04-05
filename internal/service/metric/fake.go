package metric

import (
	"context"
	"time"

	"github.com/slok/meterm/internal/model"
)

// FakeGatherer is a fake Gatherer.
type FakeGatherer struct{}

// GatherSingle satisfies metric.Gatherer interface.
func (f FakeGatherer) GatherSingle(_ context.Context, _ string, t time.Time) ([]model.MetricSeries, error) {
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
