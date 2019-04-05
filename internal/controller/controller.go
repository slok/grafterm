package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/slok/meterm/internal/model"
	"github.com/slok/meterm/internal/service/metric"
)

// Controller is what has the domain logic, the one that
// can translate from the views to the models.
type Controller interface {
	// GetSingleMetric will get one single metric value at a point in time.
	GetSingleMetric(ctx context.Context, query string, t time.Time) (*model.Metric, error)
	// GetSingleInstantMetric will get one single metric value in real time.
	GetSingleInstantMetric(ctx context.Context, query string) (*model.Metric, error)
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

func (c controller) GetSingleMetric(ctx context.Context, query string, t time.Time) (*model.Metric, error) {
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

func (c controller) GetSingleInstantMetric(ctx context.Context, query string) (*model.Metric, error) {
	return c.GetSingleMetric(ctx, query, time.Now().UTC())
}
