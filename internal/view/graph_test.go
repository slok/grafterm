package view_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mcontroller "github.com/slok/meterm/internal/mocks/controller"
	mrender "github.com/slok/meterm/internal/mocks/view/render"
	"github.com/slok/meterm/internal/model"
	"github.com/slok/meterm/internal/service/log"
	"github.com/slok/meterm/internal/view"
	"github.com/slok/meterm/internal/view/render"
)

// helper function to convert a float to a render.Value pointer.
func rv(f float64) *render.Value {
	v := render.Value(f)
	return &v
}

func TestGraphWidget(t *testing.T) {
	// Common precalculated data that should be expected.
	t1, _ := time.Parse(time.RFC3339, "2019-04-13T09:30:00+00:00")
	t1Minus100m := t1.Add(-100 * time.Minute)
	xLabels := []string{"07:50:00", "08:00:00", "08:10:00", "08:20:00", "08:30:00", "08:40:00", "08:50:00", "09:00:00", "09:10:00", "09:20:00"}
	graphCapacity := 10

	tests := []struct {
		name   string
		appCfg view.AppConfig
		cfg    model.Widget
		exp    func(*testing.T, *mcontroller.Controller, *mrender.GraphWidget)
		expErr bool
	}{
		{
			name: "A graph without without capacity on the terminal should no render anything.",
			appCfg: view.AppConfig{
				RefreshInterval: 1 * time.Second,
			},
			cfg: model.Widget{
				WidgetSource: model.WidgetSource{
					Graph: &model.GraphWidgetSource{},
				},
			},
			exp: func(t *testing.T, mc *mcontroller.Controller, mg *mrender.GraphWidget) {
				mg.On("GetGraphPointQuantity").Return(0)
			},
		},
		{
			name: "A graph with all data points should render all values.",
			appCfg: view.AppConfig{
				RefreshInterval: 1 * time.Second,
				TimeRangeEnd:    t1,
				TimeRangeStart:  t1Minus100m,
			},
			cfg: model.Widget{
				WidgetSource: model.WidgetSource{
					Graph: &model.GraphWidgetSource{
						Queries: []model.Query{
							model.Query{Expr: "test"},
						},
					},
				},
			},
			exp: func(t *testing.T, mc *mcontroller.Controller, mg *mrender.GraphWidget) {
				mg.On("GetGraphPointQuantity").Return(graphCapacity)

				// Having all datapoints we should render all the values.
				seriess := []model.MetricSeries{
					model.MetricSeries{
						ID: "test",
						Metrics: []model.Metric{
							model.Metric{Value: 1, TS: t1Minus100m.Add(1 * time.Minute)},
							model.Metric{Value: 2, TS: t1Minus100m.Add(12 * time.Minute)},
							model.Metric{Value: 3, TS: t1Minus100m.Add(21 * time.Minute)},
							model.Metric{Value: 4, TS: t1Minus100m.Add(39 * time.Minute)},
							model.Metric{Value: 5, TS: t1Minus100m.Add(46 * time.Minute)},
							model.Metric{Value: 6, TS: t1Minus100m.Add(53 * time.Minute)},
							model.Metric{Value: 7, TS: t1Minus100m.Add(66 * time.Minute)},
							model.Metric{Value: 8, TS: t1Minus100m.Add(71 * time.Minute)},
							model.Metric{Value: 9, TS: t1Minus100m.Add(85 * time.Minute)},
							model.Metric{Value: 10, TS: t1Minus100m.Add(92 * time.Minute)},
						},
					},
				}

				// Check it gets the step correctly.
				expStep := 10 * time.Minute
				mc.On("GetRangeMetrics", mock.Anything, mock.Anything, t1Minus100m, t1, expStep).Return(seriess, nil)

				// Check the data for rendering is correctly calculated.
				// Buckets index based on time (check xLabels to a fast view).
				values := []*render.Value{rv(1), rv(2), rv(3), rv(4), rv(5), rv(6), rv(7), rv(8), rv(9), rv(10)}

				series := []render.Series{
					render.Series{
						Label:   "test",
						Color:   "#7EB26D", // First color.
						XLabels: xLabels,
						Values:  values,
					},
				}
				mg.On("Sync", series).Return(nil)
			},
		},
		{
			name: "A graph with no data points at the begginning should ignore these first values.",
			appCfg: view.AppConfig{
				RefreshInterval: 1 * time.Second,
				TimeRangeEnd:    t1,
				TimeRangeStart:  t1Minus100m,
			},
			cfg: model.Widget{
				WidgetSource: model.WidgetSource{
					Graph: &model.GraphWidgetSource{
						Queries: []model.Query{
							model.Query{Expr: "test"},
						},
					},
				},
			},
			exp: func(t *testing.T, mc *mcontroller.Controller, mg *mrender.GraphWidget) {
				mg.On("GetGraphPointQuantity").Return(graphCapacity)

				// Having all datapoints we should render all the values.
				seriess := []model.MetricSeries{
					model.MetricSeries{
						ID: "test",
						Metrics: []model.Metric{
							model.Metric{Value: 5, TS: t1Minus100m.Add(46 * time.Minute)},
							model.Metric{Value: 6, TS: t1Minus100m.Add(53 * time.Minute)},
							model.Metric{Value: 7, TS: t1Minus100m.Add(66 * time.Minute)},
							model.Metric{Value: 8, TS: t1Minus100m.Add(71 * time.Minute)},
							model.Metric{Value: 9, TS: t1Minus100m.Add(85 * time.Minute)},
							model.Metric{Value: 10, TS: t1Minus100m.Add(92 * time.Minute)},
						},
					},
				}

				// Check it gets the step correctly.
				expStep := 10 * time.Minute
				mc.On("GetRangeMetrics", mock.Anything, mock.Anything, t1Minus100m, t1, expStep).Return(seriess, nil)

				// Check the data for rendering is correctly calculated.
				// Buckets index based on time (check xLabels to a fast view).
				values := []*render.Value{nil, nil, nil, nil, rv(5), rv(6), rv(7), rv(8), rv(9), rv(10)}

				series := []render.Series{
					render.Series{
						Label:   "test",
						Color:   "#7EB26D", // First color.
						XLabels: xLabels,
						Values:  values,
					},
				}
				mg.On("Sync", series).Return(nil)
			},
		},
		{
			name: "A graph with no data points at in-between should ignore these values and make no value hops on the resulting values.",
			appCfg: view.AppConfig{
				RefreshInterval: 1 * time.Second,
				TimeRangeEnd:    t1,
				TimeRangeStart:  t1Minus100m,
			},
			cfg: model.Widget{
				WidgetSource: model.WidgetSource{
					Graph: &model.GraphWidgetSource{
						Queries: []model.Query{
							model.Query{Expr: "test"},
						},
					},
				},
			},
			exp: func(t *testing.T, mc *mcontroller.Controller, mg *mrender.GraphWidget) {
				mg.On("GetGraphPointQuantity").Return(graphCapacity)

				// Having all datapoints we should render all the values.
				seriess := []model.MetricSeries{
					model.MetricSeries{
						ID: "test",
						Metrics: []model.Metric{
							model.Metric{Value: 5, TS: t1Minus100m.Add(46 * time.Minute)},
							model.Metric{Value: 6, TS: t1Minus100m.Add(53 * time.Minute)},
							model.Metric{Value: 9, TS: t1Minus100m.Add(85 * time.Minute)},
							model.Metric{Value: 10, TS: t1Minus100m.Add(92 * time.Minute)},
						},
					},
				}

				// Check it gets the step correctly.
				expStep := 10 * time.Minute
				mc.On("GetRangeMetrics", mock.Anything, mock.Anything, t1Minus100m, t1, expStep).Return(seriess, nil)

				// Check the data for rendering is correctly calculated.
				// Buckets index based on time (check xLabels to a fast view).
				values := []*render.Value{nil, nil, nil, nil, rv(5), rv(6), nil, nil, rv(9), rv(10)}

				series := []render.Series{
					render.Series{
						Label:   "test",
						Color:   "#7EB26D", // First color.
						XLabels: xLabels,
						Values:  values,
					},
				}
				mg.On("Sync", series).Return(nil)
			},
		},
		{
			name: "A graph with no data points at the end should ignore these values and make no values on the end values.",
			appCfg: view.AppConfig{
				RefreshInterval: 1 * time.Second,
				TimeRangeEnd:    t1,
				TimeRangeStart:  t1Minus100m,
			},
			cfg: model.Widget{
				WidgetSource: model.WidgetSource{
					Graph: &model.GraphWidgetSource{
						Queries: []model.Query{
							model.Query{Expr: "test"},
						},
					},
				},
			},
			exp: func(t *testing.T, mc *mcontroller.Controller, mg *mrender.GraphWidget) {
				mg.On("GetGraphPointQuantity").Return(graphCapacity)

				// Having all datapoints we should render all the values.
				seriess := []model.MetricSeries{
					model.MetricSeries{
						ID: "test",
						Metrics: []model.Metric{
							model.Metric{Value: 5, TS: t1Minus100m.Add(46 * time.Minute)},
							model.Metric{Value: 6, TS: t1Minus100m.Add(53 * time.Minute)},
							model.Metric{Value: 7, TS: t1Minus100m.Add(66 * time.Minute)},
							model.Metric{Value: 8, TS: t1Minus100m.Add(71 * time.Minute)},
						},
					},
				}

				// Check it gets the step correctly.
				expStep := 10 * time.Minute
				mc.On("GetRangeMetrics", mock.Anything, mock.Anything, t1Minus100m, t1, expStep).Return(seriess, nil)

				// Check the data for rendering is correctly calculated.
				// Buckets index based on time (check xLabels to a fast view).
				values := []*render.Value{nil, nil, nil, nil, rv(5), rv(6), rv(7), rv(8), nil, nil}

				series := []render.Series{
					render.Series{
						Label:   "test",
						Color:   "#7EB26D", // First color.
						XLabels: xLabels,
						Values:  values,
					},
				}
				mg.On("Sync", series).Return(nil)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			// Mocks.
			mgraph := &mrender.GraphWidget{}
			mgraph.On("GetWidgetCfg").Once().Return(test.cfg)
			mc := &mcontroller.Controller{}
			test.exp(t, mc, mgraph)
			mr := &mrender.Renderer{}
			mr.On("LoadDashboard", mock.Anything, mock.Anything).Once().Return([]render.Widget{mgraph}, nil)

			var err error
			ctx, cancel := context.WithCancel(context.Background())
			go func() {
				app := view.NewApp(test.appCfg, mc, mr, log.Dummy)
				err = app.Run(ctx, model.Dashboard{})
			}()

			// Give time to sync.
			time.Sleep(10 * time.Millisecond)
			cancel()

			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				mr.AssertExpectations(t)
				mc.AssertExpectations(t)
				mgraph.AssertExpectations(t)
			}
		})
	}
}
