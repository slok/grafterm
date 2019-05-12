package view

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/slok/grafterm/internal/controller"
	"github.com/slok/grafterm/internal/model"
	"github.com/slok/grafterm/internal/service/log"
	"github.com/slok/grafterm/internal/service/unit"
	"github.com/slok/grafterm/internal/view/grid"
	"github.com/slok/grafterm/internal/view/render"
	"github.com/slok/grafterm/internal/view/template"
	"github.com/slok/grafterm/internal/view/variable"
)

// AppConfig are the options to run the app.
// this configuration  has values at global app level.
type AppConfig struct {
	RefreshInterval   time.Duration
	TimeRangeStart    time.Time // Fixed optional time.
	TimeRangeEnd      time.Time // Fixed optional time.
	RelativeTimeRange time.Duration
	// OverrideVariables have priority over all variables and will override the
	// loaded variables in the dashboard.
	OverrideVariables map[string]string
}

func (a *AppConfig) defaults() {
	const (
		defRelativeTimeRange = 1 * time.Hour
		defRefreshInterval   = 10 * time.Second
	)

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
	variablers map[string]variable.Variabler

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

	// Create variablers.
	vs, err := variable.NewVariablers(variable.FactoryConfig{
		TimeRange: a.cfg.RelativeTimeRange,
		Dashboard: dashboard,
	})
	if err != nil {
		return err
	}
	a.variablers = vs

	a.running = true
	// TODO(slok): Think if we should set running to false, for now we
	// don't want to reuse the app.
	return a.run(ctx, dashboard)
}

func (a *App) run(ctx context.Context, dashboard model.Dashboard) error {
	// Create grid to render.
	var gr *grid.Grid
	var err error
	if dashboard.Grid.FixedWidgets {
		gr, err = grid.NewFixedGrid(dashboard.Grid.MaxWidth, dashboard.Widgets)
		if err != nil {
			return err
		}
	} else {
		gr, err = grid.NewAdaptiveGrid(dashboard.Grid.MaxWidth, dashboard.Widgets)
		if err != nil {
			return err
		}
	}

	// Call the View to load the dashboard and return us the widgets that we will need to call.
	renderWidgets, err := a.renderer.LoadDashboard(ctx, gr)
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

		// Dashboard data.
		dashboardData := a.getDashboardVariableData()
		overrideData := a.getOverrideVariableData()

		// Widget middlewares.
		w = withWidgetDataMiddleware(dashboardData, overrideData, w) // Assign static data to widget.

		widgets = append(widgets, w)
	}

	return widgets
}

func (a *App) getSyncConfig() syncConfig {
	cfg := syncConfig{
		timeRangeStart: a.cfg.TimeRangeStart,
		timeRangeEnd:   a.cfg.TimeRangeEnd,
	}

	// If we don't have fixed time, make the time ranges work in relative mode
	// based on now timestamp.
	if cfg.timeRangeEnd.IsZero() {
		cfg.timeRangeEnd = time.Now().UTC()
	}
	if cfg.timeRangeStart.IsZero() {
		cfg.timeRangeStart = cfg.timeRangeEnd.Add(-1 * a.cfg.RelativeTimeRange)
	}

	// Create the template data for each sync.
	cfg.templateData = a.getSyncVariableData(cfg)

	return cfg
}

func (a *App) getDashboardVariableData() template.Data {
	data := template.Data(map[string]interface{}{
		"__range":          unit.DurationToSimpleString(a.cfg.RelativeTimeRange),
		"__refresInterval": unit.DurationToSimpleString(a.cfg.RefreshInterval),
	})

	// Load variablers data from the dashboard scope.
	dashboardData := map[string]interface{}{}
	for vid, v := range a.variablers {
		if v.Scope() == variable.ScopeDashboard {
			dashboardData[vid] = v.GetValue()
		}
	}

	// Merge them.
	data = data.WithData(dashboardData)
	return data
}

func (a *App) getOverrideVariableData() template.Data {
	od := map[string]interface{}{}
	for k, v := range a.cfg.OverrideVariables {
		od[k] = v
	}
	return template.Data(od)
}

func (a *App) getSyncVariableData(cfg syncConfig) template.Data {
	data := map[string]interface{}{
		"__start": fmt.Sprintf("%v", cfg.timeRangeStart),
		"__end":   fmt.Sprintf("%v", cfg.timeRangeEnd),
	}

	for vid, v := range a.variablers {
		if v.Scope() == variable.ScopeSync {
			data[vid] = v.GetValue()
		}
	}

	return data
}
