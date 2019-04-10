package metric

import (
	"context"
	"time"

	"github.com/slok/meterm/internal/model"
)

// Gatherer knows how to gather metrics from different backends.
type Gatherer interface {
	// GatherSingle gathers one single metric at a point in time.
	GatherSingle(ctx context.Context, query model.Query, t time.Time) ([]model.MetricSeries, error)
}
