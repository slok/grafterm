package render

import (
	"context"

	"github.com/slok/meterm/internal/model"
)

// Renderer is the interface that knows how to load a dashboard to be rendered
// in some target of UI.
type Renderer interface {
	LoadDashboard(ctx context.Context, dashboard model.Dashboard) ([]Widget, error)
	Close()
}

// Widget represnets a widget that can be rendered on the view.
type Widget interface {
	GetWidgetCfg() model.Widget
}

// GaugeWidget knows how to render a Gauge kind widget that can be in percent
// or not and supports color changes.
type GaugeWidget interface {
	Widget
	Sync(isPercent bool, value float64) error
	SetColor(hexColor string) error
}

// SinglestatWidget knows how to render a Singlestat kind widget that can render text
// and supports changing color.
type SinglestatWidget interface {
	Widget
	Sync(value float64) error
	SetColor(hexColor string) error
}
