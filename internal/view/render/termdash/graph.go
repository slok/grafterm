package termdash

import (
	"fmt"
	"math"

	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/widgets/linechart"
	"github.com/mum4k/termdash/widgets/text"

	"github.com/slok/grafterm/internal/model"
	"github.com/slok/grafterm/internal/view/render"
)

const (
	fullPerc              = 99
	graphHorizontalPerc   = 80
	legendHorizontalPerc  = 19
	paddingHorizontalPerc = 10
	graphVerticalPerc     = 90
	legendVerticalPerc    = 4
	paddingVerticalPerc   = 50
	legendCharacter       = `⠤⠤`
	axesColor             = 8
	yAxisLabelsColor      = 15
	xAxisLabelsColor      = 248
)

// graph satisfies render.GraphWidget interface.
type graph struct {
	cfg model.Widget

	widgetGraph  *linechart.LineChart
	widgetLegend *text.Text
	element      grid.Element
}

func newGraph(cfg model.Widget) (*graph, error) {
	// Create the Graphwidget.
	// TODO(slok): Allow configuring the color of the axis.
	lc, err := linechart.New(
		linechart.AxesCellOpts(cell.FgColor(cell.ColorNumber(axesColor))),
		linechart.YLabelCellOpts(cell.FgColor(cell.ColorNumber(yAxisLabelsColor))),
		linechart.XLabelCellOpts(cell.FgColor(cell.ColorNumber(xAxisLabelsColor))),
		linechart.YAxisAdaptive(),
	)
	if err != nil {
		return nil, err
	}

	// If we don't need a legend then use only the graph.
	var element grid.Element
	var txt *text.Text
	if !cfg.Graph.Visualization.Legend.Disable {
		txt, err = text.New(text.WrapAtRunes())
		if err != nil {
			return nil, err
		}
	}

	element = elementFromGraphAndLegend(cfg, lc, txt)

	return &graph{
		widgetGraph:  lc,
		widgetLegend: txt,
		cfg:          cfg,
		element:      element,
	}, nil
}

func elementFromGraphAndLegend(cfg model.Widget, graph *linechart.LineChart, legend *text.Text) grid.Element {
	graphElement := grid.Widget(graph)

	elements := []grid.Element{}
	switch {
	// Disabled (no legend element).
	case cfg.Graph.Visualization.Legend.Disable:
		elements = []grid.Element{
			grid.ColWidthPerc(fullPerc, graphElement),
		}
	// To the right(elements composed by columns).
	case cfg.Graph.Visualization.Legend.RightSide:
		legendElement := grid.ColWidthPercWithOpts(
			fullPerc,
			[]container.Option{container.PaddingLeftPercent(paddingHorizontalPerc)},
			grid.Widget(legend))

		elements = []grid.Element{
			grid.ColWidthPerc(graphHorizontalPerc, graphElement),
			grid.ColWidthPerc(legendHorizontalPerc, legendElement),
		}
	// At the bottom(elements composed by rows).
	default:
		legendElement := grid.RowHeightPercWithOpts(
			fullPerc,
			[]container.Option{container.PaddingTopPercent(paddingVerticalPerc)},
			grid.Widget(legend))

		elements = []grid.Element{
			grid.RowHeightPerc(graphVerticalPerc, graphElement),
			grid.RowHeightPerc(legendVerticalPerc, legendElement),
		}
	}

	opts := []container.Option{
		container.Border(linestyle.Light),
		container.BorderTitle(cfg.Title),
	}
	element := grid.RowHeightPercWithOpts(fullPerc, opts, elements...)

	return element
}

func (g *graph) getElement() grid.Element {
	return g.element
}

func (g *graph) GetWidgetCfg() model.Widget {
	return g.cfg
}

func (g *graph) Sync(series []render.Series) error {
	// Reset legend on each sync.
	if !g.cfg.Graph.Visualization.Legend.Disable {
		g.widgetLegend.Reset()
	}

	for _, s := range series {
		// We fail all the graph sync if one of the series fail.
		err := g.syncSeries(s)
		if err != nil {
			return err
		}
	}
	return nil
}

// syncSeries will sync the widgets with one of the series.
func (g *graph) syncSeries(series render.Series) error {
	color, err := colorHexToTermdash(series.Color)
	if err != nil {
		return err
	}

	err = g.syncGraph(series, color)
	if err != nil {
		return err
	}

	err = g.syncLegend(series, color)
	if err != nil {
		return err
	}

	return nil
}

// syncGraph will set one series of metrics on the graph.
func (g *graph) syncGraph(series render.Series, color cell.Color) error {
	// Convert to float64 values.
	values := make([]float64, len(series.Values))
	for i, value := range series.Values {
		// Use NaN as no value for Termdash.
		v := math.NaN()
		if value != nil {
			v = float64(*value)
		}
		values[i] = v
	}

	// Sync widget.
	err := g.widgetGraph.Series(series.Label, values,
		linechart.SeriesCellOpts(cell.FgColor(color)),
		linechart.SeriesXLabels(g.xLabelsSliceToMap(series.XLabels)))
	if err != nil {
		return err
	}

	return nil
}

// syncLegend will set the legend if required and with the correct format.
func (g *graph) syncLegend(series render.Series, color cell.Color) error {
	legend := ""
	switch {
	// Disabled.
	case g.cfg.Graph.Visualization.Legend.Disable:
		return nil
	// To the right.
	case g.cfg.Graph.Visualization.Legend.RightSide:
		legend = fmt.Sprintf("%s %s\n", legendCharacter, series.Label)
	// At the bottom.
	default:
		legend = fmt.Sprintf("%s %s  ", legendCharacter, series.Label)
	}

	// Write the legend on the widget.
	err := g.widgetLegend.Write(legend, text.WriteCellOpts(cell.FgColor(color)))
	if err != nil {
		return err
	}

	return nil
}

func (g *graph) GetGraphPointQuantity() int {
	return g.widgetGraph.ValueCapacity()
}

func (g *graph) xLabelsSliceToMap(labels []string) map[int]string {
	mlabel := map[int]string{}
	for i, label := range labels {
		mlabel[i] = label
	}
	return mlabel
}
