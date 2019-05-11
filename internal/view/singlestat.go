package view

import (
	"context"
	"fmt"
	"sort"

	"github.com/slok/grafterm/internal/controller"
	"github.com/slok/grafterm/internal/model"
	"github.com/slok/grafterm/internal/service/unit"
	"github.com/slok/grafterm/internal/view/render"
	"github.com/slok/grafterm/internal/view/template"
)

const (
	valueTemplateKey = "value"
	defValueTemplate = "{{.value}}"
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
	templatedQ := s.cfg.Singlestat.Query
	templatedQ.Expr = cfg.templateData.Render(templatedQ.Expr)
	m, err := s.controller.GetSingleInstantMetric(ctx, templatedQ)
	if err != nil {
		return fmt.Errorf("error getting single instant metric: %s", err)
	}

	// Change the widget color if required.
	err = s.changeWidgetColor(m.Value)
	if err != nil {
		return err
	}

	// Update the render view value.
	text, err := s.valueToText(cfg, m.Value)
	if err != nil {
		return fmt.Errorf("error rendering value: %s", err)
	}
	err = s.rendererWidget.Sync(text)
	if err != nil {
		return fmt.Errorf("error setting value on render view widget: %s", err)
	}

	return nil
}

func (s *singlestat) changeWidgetColor(val float64) error {
	if len(s.cfg.Singlestat.Thresholds) == 0 {
		return nil
	}

	color, err := widgetColorManager{}.GetColorFromThresholds(s.cfg.Singlestat.Thresholds, val)
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

// valueToText will use a templater to get the text. The value
// obtained for the widget will be available under the described
// key.`
func (s *singlestat) valueToText(cfg syncConfig, value float64) (string, error) {
	var templateData template.Data

	// If we have a unit set transform.
	// If unit is unset and value text template neither then apply default
	// unit transformation.
	wcfg := s.cfg.Singlestat
	if wcfg.Unit != "" || (wcfg.Unit == "" && wcfg.ValueText == "") {
		f, err := unit.NewUnitFormatter(wcfg.Unit)
		if err != nil {
			return "", err
		}
		templateData = cfg.templateData.WithData(map[string]interface{}{
			valueTemplateKey: f(value, wcfg.Decimals),
		})
	} else {
		templateData = cfg.templateData.WithData(map[string]interface{}{
			valueTemplateKey: value,
		})
	}

	vTpl := s.cfg.Singlestat.ValueText
	if vTpl == "" {
		vTpl = defValueTemplate
	}

	return templateData.Render(vTpl), nil
}
