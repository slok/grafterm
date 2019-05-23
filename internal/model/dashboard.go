package model

import (
	"fmt"
	"regexp"

	"github.com/slok/grafterm/internal/service/unit"
)

// Defaults.
const (
	// defGridMaxWidth is the default grid width used when is not set.
	defGridMaxWidth = 100
)

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
	ValueRepresentation `json:",inline"`
	Query               Query       `json:"query,omitempty"`
	ValueText           string      `json:"valueText,omitempty"`
	Thresholds          []Threshold `json:"thresholds,omitempty"`
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
	YAxis          YAxis            `json:"yAxis,omitempty"`
}

// NullPointMode is how the graph should behave when there are null
// points on the graph.
type NullPointMode string

const (
	// NullPointModeAsNull is the default mode, it will not fill the null values.
	NullPointModeAsNull NullPointMode = "null"
	// NullPointModeConnected will try to connect the null values copying the nearest values.
	NullPointModeConnected NullPointMode = "connected"
	// NullPointModeAsZero will render the null values as zeroes.
	NullPointModeAsZero NullPointMode = "zero"
)

// SeriesOverride will override visualization based on
// the regex legend.
type SeriesOverride struct {
	Regex         string         `json:"regex,omitempty"`
	CompiledRegex *regexp.Regexp `json:"-"`
	Color         string         `json:"color,omitempty"`
	NullPointMode NullPointMode  `json:"nullPointMode,omitempty"`
}

// Legend controls the legend of a widget.
type Legend struct {
	Disable   bool `json:"disable,omitempty"`
	RightSide bool `json:"rightSide,omitempty"`
}

// YAxis controls the YAxis of a widget.
type YAxis struct {
	ValueRepresentation `json:",inline"`
}

// ValueRepresentation controls the representation of a value.
type ValueRepresentation struct {
	Unit     string `json:"unit,omitempty"`
	Decimals int    `json:"decimals,omitempty"`
}

// Validate validates the object model is correct.
// A correct object means that also it will autofill the
// required default attributes so the object ends in a
// valid state.
func (d *Dashboard) Validate() error {
	err := d.Grid.validate()
	if err != nil {
		return err
	}

	for _, v := range d.Variables {
		err := v.validate()
		if err != nil {
			return err
		}
	}

	// Validate individual widgets.
	for _, w := range d.Widgets {
		err := w.validate(*d)
		if err != nil {
			return err
		}
	}

	// TODO(slok): Validate all widgets as a whole (for example total of grid)
	return nil
}

func (g *Grid) validate() error {
	if g.MaxWidth <= 0 {
		g.MaxWidth = defGridMaxWidth
	}
	return nil
}

func (v Variable) validate() error {
	if v.Name == "" {
		return fmt.Errorf("variables should have a name")
	}

	// Variable type checks.
	switch {
	case v.VariableSource.Constant != nil:
		c := v.VariableSource.Constant
		if c.Value == "" {
			return fmt.Errorf("%s constant variable needs a value", v.Name)
		}
	case v.VariableSource.Interval != nil:
		i := v.VariableSource.Interval
		if i.Steps <= 0 {
			return fmt.Errorf("%s interval variable step should be > 0", v.Name)
		}
	default:
		return fmt.Errorf("%s variable is empty, it should be of a specific type", v.Name)
	}

	return nil
}

func (w Widget) validate(d Dashboard) error {
	err := w.GridPos.validate(d.Grid)
	if err != nil {
		return fmt.Errorf("error on %s widget grid position: %s", w.Title, err)
	}

	switch {
	case w.Gauge != nil:
		err := w.Gauge.validate()
		if err != nil {
			return fmt.Errorf("error on %s gauge widget: %s", w.Title, err)
		}
	case w.Singlestat != nil:
		err := w.Singlestat.validate()
		if err != nil {
			return fmt.Errorf("error on %s singlestat widget: %s", w.Title, err)
		}
	case w.Graph != nil:
		err := w.Graph.validate()
		if err != nil {
			return fmt.Errorf("error on %s graph widget: %s", w.Title, err)
		}
	}
	return nil
}

func (g GridPos) validate(gr Grid) error {
	if g.W <= 0 {
		return fmt.Errorf("widget grid position should have a width")
	}

	if gr.FixedWidgets && g.X <= 0 {
		return fmt.Errorf("widget grid position in a fixed grid should have am X position")
	}

	if gr.FixedWidgets && g.Y <= 0 {
		return fmt.Errorf("widget grid position in a fixed grid should have am Y position")
	}

	return nil
}

func (g GaugeWidgetSource) validate() error {
	err := g.Query.validate()
	if err != nil {
		return fmt.Errorf("query error on gauge widget: %s", err)
	}

	if g.PercentValue && g.Max <= g.Min {
		return fmt.Errorf("a percent based gauge max should be greater than min")
	}

	err = validateThresholds(g.Thresholds)
	if err != nil {
		return fmt.Errorf("thresholds error on gauge widget: %s", err)
	}

	return nil
}

func (s SinglestatWidgetSource) validate() error {
	err := s.Query.validate()
	if err != nil {
		return fmt.Errorf("query error on singlestat widget: %s", err)
	}

	err = s.ValueRepresentation.validate()
	if err != nil {
		return err
	}

	err = validateThresholds(s.Thresholds)
	if err != nil {
		return fmt.Errorf("thresholds error on singlestat widget: %s", err)
	}

	return nil
}

func (g GraphWidgetSource) validate() error {
	if len(g.Queries) <= 0 {
		return fmt.Errorf("graph must have at least one query")
	}

	for _, q := range g.Queries {
		err := q.validate()
		if err != nil {
			return err
		}
	}

	sos, err := validateSeriesOverride(g.Visualization.SeriesOverride)
	if err != nil {
		return fmt.Errorf("series override error on graph widget: %s", err)
	}
	g.Visualization.SeriesOverride = sos

	err = g.Visualization.YAxis.validate()
	if err != nil {
		return err
	}

	return nil
}

func (q Query) validate() error {
	if q.Expr == "" {
		return fmt.Errorf("query must have an expression")
	}

	if q.DatasourceID == "" {
		return fmt.Errorf("query must have have a datosource ID")
	}
	return nil
}

func validateThresholds(ts []Threshold) error {
	startValues := map[float64]struct{}{}
	for _, t := range ts {
		_, ok := startValues[t.StartValue]
		if ok {
			return fmt.Errorf("threshold start value settings can't be repeated in multiple thresholds")
		}

		startValues[t.StartValue] = struct{}{}
	}

	return nil
}

func (y YAxis) validate() error {
	err := y.ValueRepresentation.validate()
	if err != nil {
		return err
	}
	return nil
}

func (v ValueRepresentation) validate() error {
	_, err := unit.NewUnitFormatter(v.Unit)
	if err != nil {
		return fmt.Errorf("%s is an invalid unit", v.Unit)
	}

	return nil
}

func (s SeriesOverride) validate() error {
	if s.Regex == "" {
		return fmt.Errorf("a graph override for series should have a regex")
	}

	err := s.NullPointMode.validate()
	if err != nil {
		return err
	}

	return nil
}

func (n *NullPointMode) validate() error {
	if *n == "" {
		*n = NullPointModeAsNull
	}

	switch *n {
	case NullPointModeAsNull, NullPointModeAsZero, NullPointModeConnected:
		return nil
	default:
		return fmt.Errorf("null point mode '%s' is not a valid mode", *n)
	}
}

func validateSeriesOverride(sos []SeriesOverride) ([]SeriesOverride, error) {
	regexes := map[string]struct{}{}
	for i, s := range sos {
		err := s.validate()
		if err != nil {
			return sos, err
		}

		_, ok := regexes[s.Regex]
		if ok {
			return sos, fmt.Errorf("series override regex setting can't be repeated in multiple series override")
		}
		regexes[s.Regex] = struct{}{}

		// Compile the regex.
		re, err := regexp.Compile(s.Regex)
		if err != nil {
			return sos, err
		}
		s.CompiledRegex = re
		sos[i] = s
	}

	return sos, nil
}
