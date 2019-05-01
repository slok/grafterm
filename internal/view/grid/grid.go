package grid

import (
	"math"
	"sort"

	"github.com/slok/grafterm/internal/model"
)

const (
	maxWidthPercent = 100
)

// Element is a "placeable" element on the grid, depending on the
// implementation of the renderer will be created in one way or another.
type Element struct {
	// Percent size is the percent of the total in the horizontal axis.
	PercentSize int
	// Empty marks the element as an empty block that will not be used.
	Empty bool
	// Widget is the widget to be placed.
	Widget model.Widget
}

// Row is composed by multiple elements.
type Row struct {
	// Elements are the elements that will be placed on the row. The
	// elements of a row are horizontally placed. also known as
	// the X axis.
	Elements []*Element
	// PercentSize is the size in percentage of the total vertical axis.
	PercentSize int
}

// Grid is the grid itself, it's composed by rows that inside of the rows
// are the columns, the elements are the columns.
//
// ----------------------------------------------------
// [------element------] [--element--] [---element---]
// ----------------------------------------------------
// [element] [element] [------------element----------]
// ----------------------------------------------------
// [-element-]       [----element----]        [element]
// ----------------------------------------------------
type Grid struct {
	// Is the max size of the X axis. This is equal to a 100 percentage.
	MaxWidth int
	// Is the max size of the y axis. This is equal to a 100 percentage.
	MaxHeight int
	// Rows are the rows the grid has (inside the rows are the columns).
	// the rows are vertically placed, also know as the Y axis.
	Rows []*Row
}

// NewAdaptiveGrid returns a grid that places the widgets in the received order
// without checking its position (x, y) only using the size of the widgets.
// It will adapt the rows dinamically so the widgets that are bigger than the empty
// space on the row, will be placed in the next row an so own, on after the other
// creating new rows until all the widgets have been placed.
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
		if filledRow+e.PercentSize > maxWidthPercent {
			// If there is spare space on the row, before creating a new row
			// create an empty widget to fill the row until the end.
			if filledRow < maxWidthPercent {
				r.Elements = append(r.Elements, &Element{
					Empty:       true,
					PercentSize: maxWidthPercent - filledRow,
				})
			}

			// Next and new row.
			currentRow++
			filledRow = 0
			g.Rows = append(g.Rows, &Row{})
			r = g.Rows[currentRow]
		}

		// Add widget to row.
		filledRow += e.PercentSize
		r.Elements = append(r.Elements, e)
	}

	// Set the size of the rows, the rows have been dinamically created so until
	// we had all the rows we can't be sure what is the total of the vertical axis,
	// set the same vertical percent size of the rows to all the rows
	// (e.g 4 rows of 25% or 3 rows of 33% or 10 rows of 10% ).
	totalRows := len(g.Rows) // This is the 100%.
	for _, row := range g.Rows {
		row.PercentSize = percent(1, totalRows)
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
			// widget it means that we have a blank space.
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
		if rowFilled < maxWidthPercent {
			rowElements = append(rowElements, &Element{
				Empty:       true,
				PercentSize: maxWidthPercent - rowFilled,
			})
		}

		row.Elements = rowElements
	}
}

// initRows creates all the rows in empty state.
func (g *Grid) initRows() {
	for i := 0; i < g.MaxHeight; i++ {
		g.Rows = append(g.Rows, &Row{
			PercentSize: percent(1, g.MaxHeight),
		})
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
