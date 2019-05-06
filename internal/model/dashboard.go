package model

import "regexp"

// Dashboard represents a dashboard.
type Dashboard struct {
	Grid      Grid       `json:"grid,omitempty"`
	Variables []Variable `json:"variables,omitempty"`
	Widgets   []Widget   `json:"widgets,omitempty"`
}

// Variable is a dynamic variable that will be available through the
// dashboard.
type Variable struct {
	Name           string `json:"name,omitempty"`
	VariableSource `json:",inline"`
}

// VariableSource is the variable kind with it's data.
type VariableSource struct {
	Constant *ConstantVariableSource `json:"constant,omitempty"`
	Interval *IntervalVariableSource `json:"interval,omitempty"`
}

// ConstantVariableSource represents the constant variables.
type ConstantVariableSource struct {
	Value string `json:"value,omitempty"`
}

// IntervalVariableSource represents the interval variables.
type IntervalVariableSource struct {
	Steps int `json:"steps,omitempty"`
}

// Widget represents a widget.
type Widget struct {
	Title        string  `json:"title,omitempty"`
	GridPos      GridPos `json:"gridPos,omitempty"`
	WidgetSource `json:",inline"`
}

// Grid represents the options of the grid in the dashboard.
type Grid struct {
	// Fixed means that the grid positions (gridPos) of the widgets
	// will be fixed and need X and Y values.
	// If false it will be adaptive and will ignore X and Y values
	// and only use the size of the widget (W, width).
	FixedWidgets bool `json:"fixedWidgets,omitempty"`
	// MaxWidth is the maximum width (horizontal) the Grid will have, this will be
	// the scale for the widgets `GridPos.W`. For example a `GridPos.W: 50`
	// in a `Grid.MaxWidth: 100` would be half of the row, but in a `Grid.MaxWidth: 1000`
	// would be a 5% of the row.
	// Not setting MaxWidth or setting to 0 would fallback to default MaxWidth.
	MaxWidth int `json:"maxWidth,omitempty"`
}

// GridPos represents the grid position.
type GridPos struct {
	// X represents the position on the grid (from 0 to 100).
	X int `json:"x,omitempty"`
	// Y represents the position on the grid (from 0 to infinite,
	// where the total will be used using all the widgets Y and H).
	Y int `json:"y,omitempty"`
	// W represents the width of the widget (same unit as X).
	W int `json:"w,omitempty"`
	// TODO(slok): H represents the height of the widget (same unit as Y).
	// H int `json:"h,omitempty"`
}

// WidgetSource will tell what kind of widget is.
type WidgetSource struct {
	Singlestat *SinglestatWidgetSource `json:"singlestat,omitempty"`
	Gauge      *GaugeWidgetSource      `json:"gauge,omitempty"`
	Graph      *GraphWidgetSource      `json:"graph,omitempty"`
}

// SinglestatWidgetSource represents a simple value widget.
type SinglestatWidgetSource struct {
	Query      Query       `json:"query,omitempty"`
	ValueText  string      `json:"valueText,omitempty"`
	Thresholds []Threshold `json:"thresholds,omitempty"`
}

// GaugeWidgetSource represents a simple value widget in donut format.
type GaugeWidgetSource struct {
	Query        Query       `json:"query,omitempty"`
	PercentValue bool        `json:"percentValue,omitempty"`
	Max          int         `json:"max,omitempty"`
	Min          int         `json:"min,omitempty"`
	Thresholds   []Threshold `json:"thresholds,omitempty"`
}

// GraphWidgetSource represents a simple value widget in donut format.
type GraphWidgetSource struct {
	Queries       []Query            `json:"queries,omitempty"`
	Visualization GraphVisualization `json:"visualization,omitempty"`
}

// Query is the query that will be made to the datasource.
type Query struct {
	Expr string `json:"expr,omitempty"`
	// Legend accepts `text.template` format.
	Legend       string `json:"legend,omitempty"`
	DatasourceID string `json:"datasourceID,omitempty"`
}

// Threshold is a color threshold that is composed
// with the start value, 0 means the base or starting
// threshold.
type Threshold struct {
	StartValue float64 `json:"startValue"`
	Color      string  `json:"color"`
}

// GraphVisualization controls how the graph will visualize
// lines, colors, legend...
type GraphVisualization struct {
	SeriesOverride []SeriesOverride `json:"seriesOverride,omitempty"`
	Legend         Legend           `json:"legend,omitempty"`
}

// SeriesOverride will override visualization based on
// the regex legend.
type SeriesOverride struct {
	Regex         string `json:"regex,omitempty"`
	CompiledRegex *regexp.Regexp
	Color         string `json:"color,omitempty"`
}

// Legend controls the legend of a widget.
type Legend struct {
	Disable   bool `json:"disable,omitempty"`
	RightSide bool `json:"rightSide,omitempty"`
}
