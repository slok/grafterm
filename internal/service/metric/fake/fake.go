package fake

import (
	"context"
	"fmt"
	"time"

	"github.com/slok/grafterm/internal/model"
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

// GatherRange satisfies metric.Gatherer interface.
func (g Gatherer) GatherRange(ctx context.Context, query model.Query, start, end time.Time, step time.Duration) ([]model.MetricSeries, error) {
	// Get some series.
	series := []model.MetricSeries{}
	seriesQ := 3
	for i := 0; i < seriesQ; i++ {
		id := fmt.Sprintf("fake-%d", i)
		s := model.MetricSeries{
			ID: id,
			Labels: map[string]string{
				"fake": "true",
				"name": id,
			},
			Metrics: generateMetrics(i*11, start, end, step),
		}

		series = append(series, s)
	}

	return series, nil
}

func generateMetrics(offset int, start, end time.Time, step time.Duration) []model.Metric {
	metrics := []model.Metric{}
	for i := 1; ; i++ {
		t := start.Add(step * time.Duration(i))
		val := float64(offset + t.Second())

		// Add some noise.
		noise := t.Second() % 3
		if noise%2 == 0 {
			val = val + float64(noise)
		} else {
			val = val - float64(noise)
		}

		m := model.Metric{
			TS:    t,
			Value: val + float64(offset),
		}
		metrics = append(metrics, m)

		if t.After(end) {
			break
		}
	}

	return metrics
}
