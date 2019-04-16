package termdash

import (
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/widgets/linechart"

	"github.com/slok/meterm/internal/model"
	"github.com/slok/meterm/internal/view/render"
)

// graph satisfies render.GraphWidget interface.
type graph struct {
	cfg model.Widget

	widget  *linechart.LineChart
	element grid.Element
}

func newGraph(cfg model.Widget) (*graph, error) {
	// Create the widget.
	// TODO(slok): Allow configuring the color of the axis.
	lc, err := linechart.New(
		linechart.AxesCellOpts(cell.FgColor(cell.ColorRed)),
		linechart.YLabelCellOpts(cell.FgColor(cell.ColorGreen)),
		linechart.XLabelCellOpts(cell.FgColor(cell.ColorCyan)),
		linechart.YAxisAdaptive(),
	)
	if err != nil {
		return nil, err
	}

	// Create the element using the new widget.
	element := grid.Widget(lc,
		container.Border(linestyle.Light),
		container.BorderTitle(cfg.Title),
	)

	return &graph{
		widget:  lc,
		cfg:     cfg,
		element: element,
	}, nil
}

func (g *graph) getElement() grid.Element {
	return g.element
}

func (g *graph) GetWidgetCfg() model.Widget {
	return g.cfg
}

func (g *graph) Sync(series []render.Series) error {
	for _, s := range series {
		// We fail all the graph sync if one of the series fail.
		err := g.syncSeries(s)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *graph) syncSeries(series render.Series) error {
	color, err := colorHexToTermdash(series.Color)
	if err != nil {
		return err
	}

	// Convert to float64 values.
	values := make([]float64, len(series.Values))
	for i, value := range series.Values {
		// Termdash  doesn't support no values.
		// for now no values will be 0.
		// TODO(slok): Track the issue https://github.com/mum4k/termdash/issues/184
		v := 0.0
		if value != nil {
			v = float64(*value)
		}
		values[i] = v
	}

	err = g.widget.Series(series.Label, values,
		linechart.SeriesCellOpts(cell.FgColor(color)),
		linechart.SeriesXLabels(g.xLabelsSliceToMap(series.XLabels)))
	if err != nil {
		return err
	}

	return nil
}

func (g *graph) GetGraphPointQuantity() int {
	return g.widget.ValueCapacity()
}

func (g *graph) xLabelsSliceToMap(labels []string) map[int]string {
	mlabel := map[int]string{}
	for i, label := range labels {
		mlabel[i] = label
	}
	return mlabel
}
