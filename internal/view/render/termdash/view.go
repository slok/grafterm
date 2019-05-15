package termdash

import (
	"context"
	"time"

	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/keyboard"
	"github.com/mum4k/termdash/terminal/termbox"
	"github.com/mum4k/termdash/terminal/terminalapi"

	"github.com/slok/grafterm/internal/model"
	"github.com/slok/grafterm/internal/service/log"
	graftermgrid "github.com/slok/grafterm/internal/view/grid"
	"github.com/slok/grafterm/internal/view/render"
)

const (
	rootID         = "root"
	redrawInterval = 250 * time.Millisecond
)

// elementer is an internal interface that all widgets from the termdash
// render engine implementation need to implement, this way the widgets
// can create subelements by their own and the `termDashboard` does not
// to be aware, so a widget can be composed of 2 widgets under the hoods.
type elementer interface {
	getElement() grid.Element
}

// View is what renders the metrics.
type termDashboard struct {
	widgets []render.Widget
	logger  log.Logger
	cancel  func()

	// Term fields.
	terminal *termbox.Terminal
}

// NewTermDashboard returns a new terminal view, it accepts a cancel function that will
// be called when the terminal rendered quit function is called. This is required because
// the events now are captured by the rendered terminal.
func NewTermDashboard(cancel func(), logger log.Logger) (render.Renderer, error) {
	t, err := termbox.New()
	if err != nil {
		return nil, err
	}

	return &termDashboard{
		cancel:   cancel,
		terminal: t,
		logger:   logger,
	}, nil
}

func (t *termDashboard) Close() {
	t.terminal.Close()
}

// Run will run the view, its' a blocker.
func (t *termDashboard) LoadDashboard(ctx context.Context, gr *graftermgrid.Grid) ([]render.Widget, error) {
	// Create main view (root).
	c, err := container.New(t.terminal, container.ID(rootID))
	if err != nil {
		return nil, err
	}

	// Get the layout from the grid.
	gridOpts, err := t.gridLayout(gr)
	if err != nil {
		return []render.Widget{}, err
	}

	err = c.Update(rootID, gridOpts...)
	if err != nil {
		return []render.Widget{}, err
	}

	go func() {
		quitter := func(k *terminalapi.Keyboard) {
			if k.Key == 'q' || k.Key == 'Q' || k.Key == keyboard.KeyEsc {
				t.cancel()
			}
		}
		if err := termdash.Run(ctx, t.terminal, c, termdash.KeyboardSubscriber(quitter), termdash.RedrawInterval(redrawInterval)); err != nil {
			t.logger.Errorf("error running termdash terminal: %s", err)
			// TODO(slok): exit on error.
		}
	}()

	return t.widgets, nil
}

func (t *termDashboard) gridLayout(gr *graftermgrid.Grid) ([]container.Option, error) {
	builder := grid.New()

	// Create the rendering widgets.
	rowsElements := [][]grid.Element{}
	for _, row := range gr.Rows {
		rowElements := []grid.Element{}
		totalFilled := 0
		for _, rowElement := range row.Elements {
			cfg := rowElement.Widget

			// New widget.
			var element grid.Element
			if !rowElement.Empty {
				widget, err := t.newWidget(cfg)
				if err != nil {
					t.logger.Errorf("error creating widget: %s", err)
					continue
				}
				// Add widget to the tracked widgets so the app can control them.
				t.widgets = append(t.widgets, widget)

				// Get the grid.Element from our widget and place on the grid.
				element = widget.(elementer).getElement()
			}

			// Fix the size on the last element.
			// Termdash does not allow a column greater than 99, we have
			// used percents (0-100), so we remove a 1% from the last element.
			// Ugly but makes easy to work with % and is difficult for the
			// eye to notice of the 1%.
			elementPerc := rowElement.PercentSize
			if totalFilled+elementPerc >= 100 {
				elementPerc--
			}
			totalFilled += elementPerc

			// Place it on the row.
			element = grid.ColWidthPerc(elementPerc, element)
			rowElements = append(rowElements, element)
		}
		rowsElements = append(rowsElements, rowElements)
	}

	// Add rows to grid.
	var gridElements []grid.Element
	totalFilled := 0
	for i, row := range gr.Rows {
		rowElements := rowsElements[i]
		rowPerc := row.PercentSize
		// Fix the size on the last element.
		// Termdash does not allow a rows greater than 99, we have
		// used percents (0-100), so we remove a 1% from the last element.
		// Ugly but makes easy to work with % and is difficult for the
		// eye to notice of the 1%.
		if totalFilled+rowPerc >= 100 {
			rowPerc--
		}
		totalFilled += rowPerc

		// Place the row.
		rowElement := grid.RowHeightPerc(rowPerc, rowElements...)
		gridElements = append(gridElements, rowElement)
	}

	// Add rows.
	builder.Add(gridElements...)

	// Get the layout from the grid.
	return builder.Build()
}

func (t *termDashboard) newWidget(widgetcfg model.Widget) (render.Widget, error) {
	var widget render.Widget
	var err error

	switch {
	case widgetcfg.Gauge != nil:
		widget, err = newGauge(widgetcfg)
	case widgetcfg.Singlestat != nil:
		widget, err = newSinglestat(widgetcfg)
	case widgetcfg.Graph != nil:
		widget, err = newGraph(widgetcfg)
	}

	return widget, err
}
