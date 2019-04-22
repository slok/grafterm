package render

import (
	"context"

	"github.com/slok/grafterm/internal/model"
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

// Value is the value of a metric.
type Value float64

// Series are the series that can be rendered.
type Series struct {
	Label string
	Color string
	// XLabels are the labels that will be displayed on the X axis
	// the position of the label is the index of the slice.
	XLabels []string
	// Value slice, if there is no value we will use a nil value
	// we could use NaN floats but nil is more idiomatic and easy
	// to understand.
	Values []*Value
}

// GraphWidget knows how to render a Graph kind widget that renders lines in
// a two axis space using lines, dots... depending on the render implementation.
type GraphWidget interface {
	Widget
	// GetGraphPointQuantity will return the number of points the graph can display
	// on the X axis at this given moment (is a best effort, when updating the graph
	// could have changed the size).
	GetGraphPointQuantity() int
	// Sync will sync the different series on the graph.
	Sync(series []Series) error
}
