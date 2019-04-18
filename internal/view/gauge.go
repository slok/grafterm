package view

import (
	"context"
	"fmt"
	"sort"

	"github.com/slok/meterm/internal/controller"
	"github.com/slok/meterm/internal/model"
	"github.com/slok/meterm/internal/view/render"
)

// gauge is a widget that represents a metric in percent format.
type gauge struct {
	controller     controller.Controller
	rendererWidget render.GaugeWidget
	cfg            model.Widget
	currentColor   string
	syncLock       syncingFlag
}

func newGauge(controller controller.Controller, rendererWidget render.GaugeWidget) widget {
	cfg := rendererWidget.GetWidgetCfg()

	// Sort gauge thresholds. Optimization so we don't have to sort every time we calculate
	// a color.
	sort.Slice(cfg.Gauge.Thresholds, func(i, j int) bool {
		return cfg.Gauge.Thresholds[i].StartValue < cfg.Gauge.Thresholds[j].StartValue
	})

	return &gauge{
		controller:     controller,
		rendererWidget: rendererWidget,
		cfg:            cfg,
	}
}

func (g *gauge) sync(ctx context.Context, cfg syncConfig) error {
	// If already syncinc ignore call.
	if g.syncLock.Get() {
		return nil
	}
	// If didn't changed the value means some other sync process
	// already entered before us.
	if !g.syncLock.Set(true) {
		return nil
	}
	defer g.syncLock.Set(false)

	// Gather the gauge value.
	m, err := g.controller.GetSingleInstantMetric(ctx, g.cfg.Gauge.Query)
	if err != nil {
		return fmt.Errorf("error getting single instant metric: %s", err)
	}

	// calculate percent value if required.
	val := m.Value
	if g.cfg.Gauge.PercentValue {
		val = g.getPercentValue(val)
	}

	// Change the widget color if required.
	err = g.changeWidgetColor(val)
	if err != nil {
		return err
	}

	// Update the render view value.
	err = g.rendererWidget.Sync(g.cfg.Gauge.PercentValue, val)
	if err != nil {
		return fmt.Errorf("error setting value on render view widget: %s", err)
	}

	return nil
}

func (g *gauge) getPercentValue(val float64) float64 {
	// Calculate percent, if not max assume is from 0 to 100.
	if g.cfg.Gauge.Max != 0 {
		val = val - float64(g.cfg.Gauge.Min)
		cap := g.cfg.Gauge.Max - g.cfg.Gauge.Min
		val = val / float64(cap) * 100
	}

	if val > 100 {
		val = 100
	}

	if val < 0 {
		val = 0
	}

	return val
}

func (g *gauge) changeWidgetColor(val float64) error {
	if len(g.cfg.Gauge.Thresholds) == 0 {
		return nil
	}

	color, err := widgetColorManager{}.GetColorFromThresholds(g.cfg.Gauge.Thresholds, val)
	if err != nil {
		return fmt.Errorf("error getting threshold color: %s", err)
	}

	// If is the same color then don't change the widget color.
	if color == g.currentColor {
		return nil
	}

	// Change the color of the gauge widget.
	err = g.rendererWidget.SetColor(color)
	if err != nil {
		return fmt.Errorf("error setting color on view widget: %s", err)
	}

	// Update state.
	g.currentColor = color

	return nil
}
