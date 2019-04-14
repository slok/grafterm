package view

import (
	"context"
	"fmt"
	"sort"

	"github.com/slok/meterm/internal/controller"
	"github.com/slok/meterm/internal/model"
	"github.com/slok/meterm/internal/view/render"
)

// singlestat is a widget that represents in text mode.
type singlestat struct {
	controller     controller.Controller
	rendererWidget render.SinglestatWidget
	currentColor   string
	cfg            model.Widget
	syncLock       syncingFlag
}

func newSinglestat(controller controller.Controller, rendererWidget render.SinglestatWidget) widget {
	cfg := rendererWidget.GetWidgetCfg()

	// Sort widget thresholds. Optimization so we don't have to sort every time we calculate
	// a color.
	sort.Slice(cfg.Singlestat.Thresholds, func(i, j int) bool {
		return cfg.Singlestat.Thresholds[i].StartValue < cfg.Singlestat.Thresholds[j].StartValue
	})

	return &singlestat{
		controller:     controller,
		rendererWidget: rendererWidget,
		cfg:            cfg,
	}
}

func (s *singlestat) sync(ctx context.Context, cfg syncConfig) error {
	// If already syncinc ignore call.
	if s.syncLock.Get() {
		return nil
	}
	// If didn't changed the value means some other sync process
	// already entered before us.
	if !s.syncLock.Set(true) {
		return nil
	}
	defer s.syncLock.Set(false)

	// Gather the value.
	m, err := s.controller.GetSingleInstantMetric(ctx, s.cfg.Singlestat.Query)
	if err != nil {
		return fmt.Errorf("error getting single instant metric: %s", err)
	}

	// Change the widget color if required.
	err = s.changeWidgetColor(m.Value)
	if err != nil {
		return err
	}

	// Update the render view value.
	err = s.rendererWidget.Sync(m.Value)
	if err != nil {
		return fmt.Errorf("error setting value on render view widget: %s", err)
	}

	return nil

}

func (s *singlestat) changeWidgetColor(val float64) error {
	if len(s.cfg.Singlestat.Thresholds) == 0 {
		return nil
	}

	color, err := getThresholdColor(s.cfg.Singlestat.Thresholds, val)
	if err != nil {
		return fmt.Errorf("error getting threshold color: %s", err)
	}

	// If is the same color then don't change the widget color.
	if color == s.currentColor {
		return nil
	}

	// Change the color of the gauge widget.
	err = s.rendererWidget.SetColor(color)
	if err != nil {
		return fmt.Errorf("error setting color on view widget: %s", err)
	}

	// Update state.
	s.currentColor = color

	return nil
}
