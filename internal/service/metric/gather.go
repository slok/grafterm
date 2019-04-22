package metric

import (
	"context"
	"time"

	"github.com/slok/grafterm/internal/model"
)

// Gatherer knows how to gather metrics from different backends.
type Gatherer interface {
	// GatherSingle gathers one single metric at a point in time.
	GatherSingle(ctx context.Context, query model.Query, t time.Time) ([]model.MetricSeries, error)
	// GatherRange gathers multiple metrics based on a start and an end using a step duration
	// to know how many metrics needs to gather.
	// The returned metrics on the series should be ordered.
	GatherRange(ctx context.Context, query model.Query, start, end time.Time, step time.Duration) ([]model.MetricSeries, error)
}
