package model

// Dashboard represents a dashboard.
type Dashboard struct {
	Rows []Row `json:"rows,omitempty"`
}

// Row represents a row.
type Row struct {
	Title   string   `json:"title,omitempty"`
	Border  bool     `json:"border,omitempty"`
	Widgets []Widget `json:"widgets,omitempty"`
}

// Widget represents a widget.
type Widget struct {
	Title        string `json:"title,omitempty"`
	WidgetSource `json:",inline"`
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
	TextFormat string      `json:"textFormat,omitempty"`
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
	Legend  Legend  `json:"legend,omitempty"`
	Queries []Query `json:"queries,omitempty"`
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

// Legend controls the legend of a widget.
type Legend struct {
	Disable   bool `json:"disable,omitempty"`
	RightSide bool `json:"rightSide,omitempty"`
}
