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

func TestGaugeWidget(t *testing.T) {
	tests := []struct {
		name             string
		cfg              model.Widget
		controllerMetric *model.Metric
		exp              func(*mrender.GaugeWidget)
		expErr           bool
	}{
		{
			name: "A gauge without thresholds and in absolute value should render ok.",
			controllerMetric: &model.Metric{
				Value: 19,
			},
			cfg: model.Widget{
				WidgetSource: model.WidgetSource{
					Gauge: &model.GaugeWidgetSource{},
				},
			},
			exp: func(mc *mrender.GaugeWidget) {
				mc.On("Sync", false, float64(19)).Return(nil)
			},
		},
		{
			name: "A gauge without thresholds and in percent value should render ok.",
			controllerMetric: &model.Metric{
				Value: 19,
			},
			cfg: model.Widget{
				WidgetSource: model.WidgetSource{
					Gauge: &model.GaugeWidgetSource{
						PercentValue: true,
					},
				},
			},
			exp: func(mc *mrender.GaugeWidget) {
				mc.On("Sync", true, float64(19)).Return(nil)
			},
		},
		{
			name: "A gauge without thresholds and in percent value with Max and Min and Min should render ok.",
			controllerMetric: &model.Metric{
				Value: 150,
			},
			cfg: model.Widget{
				WidgetSource: model.WidgetSource{
					Gauge: &model.GaugeWidgetSource{
						PercentValue: true,
						Max:          300,
						Min:          100,
					},
				},
			},
			exp: func(mc *mrender.GaugeWidget) {
				mc.On("Sync", true, float64(25)).Return(nil)
			},
		},
		{
			name: "A gauge with (unordered) thresholds and in absolute value should set the color ok.",
			controllerMetric: &model.Metric{
				Value: 19,
			},
			cfg: model.Widget{
				WidgetSource: model.WidgetSource{
					Gauge: &model.GaugeWidgetSource{
						Thresholds: []model.Threshold{
							{Color: "#000010", StartValue: 10},
							{Color: "#000020", StartValue: 20},
							{Color: "#000005", StartValue: 5},
							{Color: "#000015", StartValue: 15},
						},
					},
				},
			},
			exp: func(mc *mrender.GaugeWidget) {
				mc.On("Sync", false, float64(19)).Return(nil)
				mc.On("SetColor", "#000015").Return(nil)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			// Mocks.
			mgauge := &mrender.GaugeWidget{}
			mgauge.On("GetWidgetCfg").Once().Return(test.cfg)
			test.exp(mgauge)

			mc := &mcontroller.Controller{}
			mc.On("GetSingleInstantMetric", mock.Anything, test.cfg.Gauge.Query).Return(test.controllerMetric, nil)
			mr := &mrender.Renderer{}
			mr.On("LoadDashboard", mock.Anything, mock.Anything).Once().Return([]render.Widget{mgauge}, nil)

			var err error
			ctx, cancel := context.WithCancel(context.Background())
			go func() {
				app := view.NewApp(view.AppConfig{}, mc, mr, log.Dummy)
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
				mgauge.AssertExpectations(t)
			}
		})
	}
}
