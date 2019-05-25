package widget_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mcontroller "github.com/slok/grafterm/internal/mocks/controller"
	mrender "github.com/slok/grafterm/internal/mocks/view/render"
	"github.com/slok/grafterm/internal/model"
	"github.com/slok/grafterm/internal/view/page/widget"
	"github.com/slok/grafterm/internal/view/sync"
	"github.com/slok/grafterm/internal/view/template"
)

func TestGaugeWidget(t *testing.T) {
	tests := []struct {
		name             string
		cfg              model.Widget
		controllerMetric *model.Metric
		syncReq          *sync.Request
		expQuery         model.Query
		exp              func(*mrender.GaugeWidget)
		expErr           bool
	}{
		{
			name: "A gauge without thresholds and in absolute value should render ok.",
			controllerMetric: &model.Metric{
				Value: 19,
			},
			syncReq: &sync.Request{},
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
			name: "A gauge should make templated queries.",
			controllerMetric: &model.Metric{
				Value: 19,
			},
			syncReq: &sync.Request{
				TemplateData: template.Data(map[string]interface{}{
					"testInterval": "10m",
				}),
			},
			cfg: model.Widget{
				WidgetSource: model.WidgetSource{
					Gauge: &model.GaugeWidgetSource{
						Query: model.Query{
							Expr: "this_is_a_test[{{ .testInterval }}]",
						},
					},
				},
			},
			expQuery: model.Query{
				Expr: "this_is_a_test[10m]",
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
			syncReq: &sync.Request{},
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
			syncReq: &sync.Request{},
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
			syncReq: &sync.Request{},
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
			mc.On("GetSingleMetric", mock.Anything, test.expQuery, mock.Anything).Return(test.controllerMetric, nil)

			var err error
			gauge := widget.NewGauge(mc, mgauge)
			gauge.Sync(context.Background(), test.syncReq)

			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				mc.AssertExpectations(t)
				mgauge.AssertExpectations(t)
			}
		})
	}
}
