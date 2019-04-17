package view

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/slok/meterm/internal/controller"
	"github.com/slok/meterm/internal/model"
	"github.com/slok/meterm/internal/service/log"
	"github.com/slok/meterm/internal/view/render"
	"github.com/slok/meterm/internal/view/template"
)

var (
	defRelativeTimeRange = 1 * time.Hour
	defRefreshInterval   = 10 * time.Second
)

// AppConfig are the options to run the app.
// this configuration  has values at global app level.
type AppConfig struct {
	RefreshInterval   time.Duration
	TimeRangeStart    time.Time // Fixed optional time.
	TimeRangeEnd      time.Time // Fixed optional time.
	RelativeTimeRange time.Duration
}

func (a *AppConfig) defaults() {
	if a.RefreshInterval == 0 {
		a.RefreshInterval = defRefreshInterval
	}
	if a.RelativeTimeRange == 0 {
		a.RelativeTimeRange = defRelativeTimeRange
	}
}

// App represents the application that will render the metrics dashboard.
type App struct {
	controller controller.Controller
	renderer   render.Renderer
	logger     log.Logger
	widgets    []widget
	cfg        AppConfig

	running bool
	mu      sync.Mutex
}

// NewApp Is the main application
func NewApp(cfg AppConfig, controller controller.Controller, renderer render.Renderer, logger log.Logger) *App {
	cfg.defaults()

	return &App{
		controller: controller,
		renderer:   renderer,
		logger:     logger,
		cfg:        cfg,
	}
}

// Run will start running the application.
func (a *App) Run(ctx context.Context, dashboard model.Dashboard) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.running {
		return errors.New("already running")
	}

	a.running = true
	// TODO(slok): Think if we should set running to false, for now we
	// don't want to reuse the app.
	return a.run(ctx, dashboard)
}

func (a *App) run(ctx context.Context, dashboard model.Dashboard) error {
	// Call the View to load the dashboard and return us the widgets that we will need to call.
	renderWidgets, err := a.renderer.LoadDashboard(ctx, dashboard)
	if err != nil {
		return err
	}

	// Create app widgets using the render widgets.
	a.widgets = a.createWidgets(renderWidgets)

	// Start the sync process. This operation blocks.
	a.sync(ctx)

	return nil
}

func (a *App) sync(ctx context.Context) {
	a.syncWidgets()

	tk := time.NewTicker(a.cfg.RefreshInterval)
	defer tk.Stop()
	for {
		// Check if we already done.
		select {
		case <-ctx.Done():
			return
		case <-tk.C:
		}

		a.syncWidgets()
	}
}

func (a *App) syncWidgets() {
	ctx := context.Background()
	cfg := a.getSyncConfig()

	// Sync all widgets.
	for _, w := range a.widgets {
		w := w
		go func() {
			// Don't wait to sync all at the same time, the widgets
			// should control multiple calls to sync and reject the sync
			// if already syncing.
			err := w.sync(ctx, cfg)
			if err != nil {
				a.logger.Errorf("error syncing widget: %s", err)
			}
		}()
	}
}

func (a *App) createWidgets(rws []render.Widget) []widget {
	widgets := []widget{}

	// Create app widgets based on the render view widgets.
	for _, rw := range rws {
		var w widget

		// Depending on the type create a widget kind or another.
		switch v := rw.(type) {
		case render.GaugeWidget:
			w = newGauge(a.controller, v)
		case render.SinglestatWidget:
			w = newSinglestat(a.controller, v)
		case render.GraphWidget:
			w = newGraph(a.controller, v, a.logger)
		default:
			continue
		}

		widgets = append(widgets, w)
	}

	return widgets
}

func (a *App) getSyncConfig() syncConfig {
	cfg := syncConfig{
		timeRangeStart: a.cfg.TimeRangeStart,
		timeRangeEnd:   a.cfg.TimeRangeEnd,
	}

	if cfg.timeRangeEnd.IsZero() {
		cfg.timeRangeEnd = time.Now()
	}

	if cfg.timeRangeStart.IsZero() {
		cfg.timeRangeStart = cfg.timeRangeEnd.Add(-1 * a.cfg.RelativeTimeRange)
	}

	// Create the template data for each sync.
	cfg.templateData = template.Data{
		Dashboard: template.Dashboard{
			Range: fmt.Sprintf("%v", a.cfg.RelativeTimeRange),
		},
	}

	return cfg
}
