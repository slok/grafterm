package grid

import (
	"math"
	"sort"

	"github.com/slok/grafterm/internal/model"
)

// Element is a "placeable" element on the grid, depending on the
// implementation of the renderer will be created in one way or another.
type Element struct {
	// Percent size is the percent of the total in the horizontal position.
	PercentSize int
	// Empty marks the element as an empty block that will not be used.
	Empty bool
	// Widget is the widget to be placed.
	Widget model.Widget
}

// Row is composed by multiple rows, the rows are horizontally placed,
// also know as the X axis.
type Row struct {
	// Elements are the elements that will be placed on the row.
	Elements    []*Element
	PercentSize int
}

// Grid is the grid itself, it's composed by rows that inside of the rows
// are the columns, the elements are the columns.
//
// ----------------------------------------------------
// [   element  ] [element] [element] [element]
// ----------------------------------------------------
// [element] [element] [            element           ]
// ----------------------------------------------------
// [element]             [element]            [element]
// ----------------------------------------------------
type Grid struct {
	// Is the max size of the X axis. This is equal to a 100 percentage.
	MaxWidth int
	// Is the max size of the y axis. This is equal to a 100 percentage.
	MaxHeight int
	// Rows are the rows the grid has (inside the rows are the columns).
	Rows []*Row
}

// NewAdaptiveGrid returns a grid that places the widgets in the received order
// without checking its position (x, y) only using the size of the widgets.
// it will adapt the rows if the widgets don't enter in the row there is being
// placed one after the other.
func NewAdaptiveGrid(maxWidth int, widgets []model.Widget) (*Grid, error) {
	d := &Grid{
		MaxWidth: maxWidth,
	}

	d.fillAdaptiveGrid(widgets)
	return d, nil
}

func (g *Grid) fillAdaptiveGrid(widgets []model.Widget) {
	g.Rows = []*Row{}
	filledRow := 0
	currentRow := 0
	for _, cfg := range widgets {
		// Initial row check existence.
		if len(g.Rows) <= currentRow {
			g.Rows = append(g.Rows, &Row{})
		}

		r := g.Rows[currentRow]

		// Create the widget element.
		e := &Element{
			PercentSize: percent(cfg.GridPos.W, g.MaxWidth),
			Widget:      cfg,
		}

		// To get he correct row of the widget then we need to see if
		// the widget is from this row or next row .
		// TODO(slok): check if widget is greater than grid totalX
		if filledRow+cfg.GridPos.W > g.MaxWidth {
			// If there is spare space on the row, before creating a new row
			// create an empty widget to fill the row until the end.
			if filledRow < g.MaxWidth {
				r.Elements = append(r.Elements, &Element{
					Empty:       true,
					PercentSize: percent(g.MaxWidth-filledRow, g.MaxWidth),
				})
			}

			// Next and new row.
			currentRow++
			filledRow = 0
			g.Rows = append(g.Rows, &Row{})
			r = g.Rows[currentRow]
		}

		// Add widget to row.
		filledRow += cfg.GridPos.W
		r.Elements = append(r.Elements, e)
	}

	// With all the grid filled, get each row size (all the same for now).
	// TODO(slok): Get highest H in the row to set the row size.
	totalRows := len(g.Rows)
	totalHeigh := 0
	for _, row := range g.Rows {
		row.PercentSize = percent(1, totalRows)
		totalHeigh += row.PercentSize
	}
}

// NewFixedGrid will place the widgets on the grid using the size and position
// of the widgets letting empty elements between them if required. This kind
// of grid needs the widgets to be exactly placed on the grid it doesn't adapt
// horizontally nor vertically.
func NewFixedGrid(maxWidth int, widgets []model.Widget) (*Grid, error) {
	maxHeight := 0
	for _, w := range widgets {
		if maxHeight < w.GridPos.Y+1 {
			maxHeight = w.GridPos.Y + 1
		}
	}

	g := &Grid{
		MaxHeight: maxHeight,
		MaxWidth:  maxWidth,
		Rows:      []*Row{},
	}

	g.fillFixedGrid(widgets)

	return g, nil
}

func (g *Grid) fillFixedGrid(widgets []model.Widget) {
	sortwidgets(widgets)
	g.initRows()
	// Create the widgets.
	for _, cfg := range widgets {
		row := g.Rows[cfg.GridPos.Y]
		row.Elements = append(row.Elements, &Element{
			PercentSize: percent(cfg.GridPos.W, g.MaxWidth),
			Widget:      cfg,
		})
	}

	// Fill the blank spaces between widgets for each row.
	for _, row := range g.Rows {
		rowFilled := 0
		var rowElements []*Element
		for _, rowElement := range row.Elements {
			posperc := percent(rowElement.Widget.GridPos.X, g.MaxWidth)

			// If what we filled is not the start point of the current
			// widget it meas that we have a black space.
			if rowFilled < posperc {
				rowElements = append(rowElements, &Element{
					Empty:       true,
					PercentSize: posperc - rowFilled,
				})
			}
			rowElements = append(rowElements, rowElement)
			rowFilled = posperc + rowElement.PercentSize
		}

		// Check if we need to fill with blank space until the end
		// of the row.
		if rowFilled < 100 {
			rowElements = append(rowElements, &Element{
				Empty:       true,
				PercentSize: 100 - rowFilled,
			})
		}

		row.Elements = rowElements
	}
}

// initRows creates all the rows in empty state.
func (g *Grid) initRows() {
	filled := 0
	for i := 0; i < g.MaxHeight; i++ {
		size := percent(1, g.MaxHeight)
		g.Rows = append(g.Rows, &Row{
			PercentSize: size,
		})
		filled += size
	}
}

func percent(value, total int) int {
	perc := float64(value) * 100 / float64(total)
	return int(math.Round(perc))
}

// sortwidgets sorts the widgets in left-right and top-down
// order.
func sortwidgets(widgets []model.Widget) {
	sort.Slice(widgets, func(i, j int) bool {
		gpi := widgets[i].GridPos
		gpj := widgets[j].GridPos

		switch {
		case gpi.Y > gpj.Y:
			return false
		case gpi.Y < gpj.Y:
			return true
		case gpi.X < gpj.X:
			return true
		case gpi.X > gpj.X:
			return false
		// If are the same then doesn't matter, we shouldn't reach here.
		default:
			return true
		}
	})
}
