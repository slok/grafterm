// Code generated by mockery v1.0.0. DO NOT EDIT.

package metric

import context "context"

import mock "github.com/stretchr/testify/mock"
import model "github.com/slok/grafterm/internal/model"
import time "time"

// Gatherer is an autogenerated mock type for the Gatherer type
type Gatherer struct {
	mock.Mock
}

// GatherRange provides a mock function with given fields: ctx, query, start, end, step
func (_m *Gatherer) GatherRange(ctx context.Context, query model.Query, start time.Time, end time.Time, step time.Duration) ([]model.MetricSeries, error) {
	ret := _m.Called(ctx, query, start, end, step)

	var r0 []model.MetricSeries
	if rf, ok := ret.Get(0).(func(context.Context, model.Query, time.Time, time.Time, time.Duration) []model.MetricSeries); ok {
		r0 = rf(ctx, query, start, end, step)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.MetricSeries)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, model.Query, time.Time, time.Time, time.Duration) error); ok {
		r1 = rf(ctx, query, start, end, step)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GatherSingle provides a mock function with given fields: ctx, query, t
func (_m *Gatherer) GatherSingle(ctx context.Context, query model.Query, t time.Time) ([]model.MetricSeries, error) {
	ret := _m.Called(ctx, query, t)

	var r0 []model.MetricSeries
	if rf, ok := ret.Get(0).(func(context.Context, model.Query, time.Time) []model.MetricSeries); ok {
		r0 = rf(ctx, query, t)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.MetricSeries)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, model.Query, time.Time) error); ok {
		r1 = rf(ctx, query, t)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
