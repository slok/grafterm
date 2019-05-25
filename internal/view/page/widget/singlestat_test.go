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

func TestSinglestatWidget(t *testing.T) {
	tests := []struct {
		name             string
		cfg              model.Widget
		syncReq          *sync.Request
		controllerMetric *model.Metric
		expQuery         model.Query
		exp              func(*mrender.SinglestatWidget)
		expErr           bool
	}{
		{
			name: "A singlestat without thresholds should render ok.",
			controllerMetric: &model.Metric{
				Value: 19.14,
			},
			syncReq: &sync.Request{},
			cfg: model.Widget{
				WidgetSource: model.WidgetSource{
					Singlestat: &model.SinglestatWidgetSource{
						ValueRepresentation: model.ValueRepresentation{
							Unit:     "none",
							Decimals: 2,
						},
					},
				},
			},
			exp: func(mc *mrender.SinglestatWidget) {
				mc.On("Sync", "19.14").Return(nil)
			},
		},
		{
			name: "A singlestat with custom template should render ok.",
			controllerMetric: &model.Metric{
				Value: 19.14,
			},
			syncReq: &sync.Request{},
			cfg: model.Widget{
				WidgetSource: model.WidgetSource{
					Singlestat: &model.SinglestatWidgetSource{
						ValueText: `this is a test with {{printf "%.1f" .value}} value`,
					},
				},
			},
			exp: func(mc *mrender.SinglestatWidget) {
				mc.On("Sync", "this is a test with 19.1 value").Return(nil)
			},
		},
		{
			name: "A singlestat should make templated queries with variables.",
			controllerMetric: &model.Metric{
				Value: 19.14,
			},
			syncReq: &sync.Request{
				TemplateData: template.Data(map[string]interface{}{
					"testInterval": "10m",
				}),
			},
			cfg: model.Widget{
				WidgetSource: model.WidgetSource{
					Singlestat: &model.SinglestatWidgetSource{
						ValueRepresentation: model.ValueRepresentation{
							Unit:     "none",
							Decimals: 2,
						},
						Query: model.Query{
							Expr: "this_is_a_test[{{ .testInterval }}]",
						},
					},
				},
			},
			expQuery: model.Query{
				Expr: "this_is_a_test[10m]",
			},
			exp: func(mc *mrender.SinglestatWidget) {
				mc.On("Sync", "19.14").Return(nil)
			},
		},
		{
			name: "A singlestat with (unordered) thresholds should set the color ok.",
			controllerMetric: &model.Metric{
				Value: 19.14,
			},
			syncReq: &sync.Request{},
			cfg: model.Widget{
				WidgetSource: model.WidgetSource{
					Singlestat: &model.SinglestatWidgetSource{
						ValueRepresentation: model.ValueRepresentation{
							Unit:     "none",
							Decimals: 2,
						},
						Thresholds: []model.Threshold{
							{Color: "#000010", StartValue: 10},
							{Color: "#000020", StartValue: 20},
							{Color: "#000005", StartValue: 5},
							{Color: "#000015", StartValue: 15},
						},
					},
				},
			},
			exp: func(mc *mrender.SinglestatWidget) {
				mc.On("Sync", "19.14").Return(nil)
				mc.On("SetColor", "#000015").Return(nil)
			},
		},
		{
			name: "A singlestat without unit should fallback to the default unit.",
			controllerMetric: &model.Metric{
				Value: 192312312321.21,
			},
			syncReq: &sync.Request{},
			cfg: model.Widget{
				WidgetSource: model.WidgetSource{
					Singlestat: &model.SinglestatWidgetSource{},
				},
			},
			exp: func(mc *mrender.SinglestatWidget) {
				mc.On("Sync", "192 Bil").Return(nil)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			// Mocks.
			msstat := &mrender.SinglestatWidget{}
			msstat.On("GetWidgetCfg").Once().Return(test.cfg)
			test.exp(msstat)

			mc := &mcontroller.Controller{}
			mc.On("GetSingleMetric", mock.Anything, test.expQuery, mock.Anything).Return(test.controllerMetric, nil)

			singlestat := widget.NewSinglestat(mc, msstat)
			err := singlestat.Sync(context.Background(), test.syncReq)

			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				mc.AssertExpectations(t)
				msstat.AssertExpectations(t)
			}
		})
	}
}
