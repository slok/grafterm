package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/slok/grafterm/internal/model"
	"github.com/slok/grafterm/internal/service/metric"
)

// Controller is what has the domain logic, the one that
// can translate from the views to the models.
type Controller interface {
	// GetSingleMetric will get one single metric value at a point in time.
	GetSingleMetric(ctx context.Context, query model.Query, t time.Time) (*model.Metric, error)
	// GetSingleInstantMetric will get one single metric value in real time.
	GetSingleInstantMetric(ctx context.Context, query model.Query) (*model.Metric, error)
	// GetRangeMetrics will get N metrics based in a time range.
	GetRangeMetrics(ctx context.Context, query model.Query, start, end time.Time, step time.Duration) ([]model.MetricSeries, error)
}

type controller struct {
	gatherer metric.Gatherer
}

// NewController returns a new controller.
func NewController(gatherer metric.Gatherer) Controller {
	return &controller{
		gatherer: gatherer,
	}
}

func (c controller) GetSingleMetric(ctx context.Context, query model.Query, t time.Time) (*model.Metric, error) {
	m, err := c.gatherer.GatherSingle(ctx, query, t)
	if err != nil {
		return nil, err
	}

	if len(m) != 1 {
		return nil, fmt.Errorf("wrong number of series returned, 1 expected, got: %d", len(m))
	}

	if len(m[0].Metrics) != 1 {
		return nil, fmt.Errorf("wrong number of metric in series returned, 1 expected, got: %d", len(m[0].Metrics))
	}

	return &m[0].Metrics[0], nil
}

func (c controller) GetSingleInstantMetric(ctx context.Context, query model.Query) (*model.Metric, error) {
	return c.GetSingleMetric(ctx, query, time.Now().UTC())
}

func (c controller) GetRangeMetrics(ctx context.Context, query model.Query, start, end time.Time, step time.Duration) ([]model.MetricSeries, error) {
	if step <= 0 {
		return nil, fmt.Errorf("step must be positive")
	}

	if !start.Before(end) {
		return nil, fmt.Errorf("start timestamp must be before end timestamp")
	}

	// Get the metrics.
	s, err := c.gatherer.GatherRange(ctx, query, start, end, step)
	if err != nil {
		return []model.MetricSeries{}, err
	}

	return s, nil
}
