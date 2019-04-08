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

func TestSinglestatWidget(t *testing.T) {
	tests := []struct {
		name             string
		cfg              model.Widget
		controllerMetric *model.Metric
		exp              func(*mrender.SinglestatWidget)
		expErr           bool
	}{
		{
			name: "A singlestat without thresholds should render ok.",
			controllerMetric: &model.Metric{
				Value: 19.14,
			},
			cfg: model.Widget{
				WidgetSource: model.WidgetSource{
					Singlestat: &model.SinglestatWidgetSource{},
				},
			},
			exp: func(mc *mrender.SinglestatWidget) {
				mc.On("Sync", 19.14).Return(nil)
			},
		},
		{
			name: "A singlestat with (unordered) thresholds should set the color ok.",
			controllerMetric: &model.Metric{
				Value: 19.14,
			},
			cfg: model.Widget{
				WidgetSource: model.WidgetSource{
					Singlestat: &model.SinglestatWidgetSource{
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
				mc.On("Sync", 19.14).Return(nil)
				mc.On("SetColor", "#000015").Return(nil)
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
			mc.On("GetSingleInstantMetric", mock.Anything, test.cfg.Singlestat.Query.Query).Return(test.controllerMetric, nil)
			mr := &mrender.Renderer{}
			mr.On("LoadDashboard", mock.Anything, mock.Anything).Once().Return([]render.Widget{msstat}, nil)

			var err error
			ctx, cancel := context.WithCancel(context.Background())
			go func() {
				app := view.NewApp(mc, mr, log.Dummy)
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
				msstat.AssertExpectations(t)
			}
		})
	}
}
