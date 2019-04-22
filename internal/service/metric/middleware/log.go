package middleware

import (
	"context"
	"time"

	"github.com/slok/grafterm/internal/model"
	"github.com/slok/grafterm/internal/service/log"
	"github.com/slok/grafterm/internal/service/metric"
)

type logger struct {
	next   metric.Gatherer
	logger log.Logger
}

// Logger is a gatherer middleware that wraps the real gatherer and logs
// the queries that it makes.
func Logger(l log.Logger, next metric.Gatherer) metric.Gatherer {
	return &logger{
		next:   next,
		logger: l,
	}
}

func (l *logger) GatherSingle(ctx context.Context, query model.Query, t time.Time) ([]model.MetricSeries, error) {
	st := time.Now()
	defer func() {
		l.logger.Infof("(%s) gathering single metric on %s: %s", time.Since(st), query.DatasourceID, query.Expr)
	}()
	return l.next.GatherSingle(ctx, query, t)
}
func (l *logger) GatherRange(ctx context.Context, query model.Query, start, end time.Time, step time.Duration) ([]model.MetricSeries, error) {
	st := time.Now()
	defer func() {
		l.logger.Infof("(%s) gathering range metric [from %v to %v with %v step] on %s: %s", time.Since(st), start, end, step, query.DatasourceID, query.Expr)
	}()
	return l.next.GatherRange(ctx, query, start, end, step)
}
