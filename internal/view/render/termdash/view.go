package termdash

import (
	"context"
	"time"

	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/termbox"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgetapi"

	"github.com/slok/meterm/internal/model"
	"github.com/slok/meterm/internal/service/log"
	"github.com/slok/meterm/internal/view/render"
)

const (
	rootID = "root"
)

// View is what renders the metrics.
type termDashboard struct {
	widgets []render.Widget
	logger  log.Logger
	cancel  func()

	// Term fields.
	terminal *termbox.Terminal
}

// NewTermDashboard returns a new terminal view, it accepts a cancel function taht will
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
func (t *termDashboard) LoadDashboard(ctx context.Context, dashboard model.Dashboard) ([]render.Widget, error) {

	// Create main view (root).
	c, err := container.New(t.terminal, container.ID(rootID))
	if err != nil {
		return nil, err
	}

	// Get the layout from the grid.
	gridOpts, err := t.gridLayout(dashboard)
	if err != nil {
		return []render.Widget{}, err
	}

	err = c.Update(rootID, gridOpts...)
	if err != nil {
		return []render.Widget{}, err
	}

	go func() {
		quitter := func(k *terminalapi.Keyboard) {
			if k.Key == 'q' || k.Key == 'Q' {
				t.cancel()
			}
		}
		if err := termdash.Run(ctx, t.terminal, c, termdash.KeyboardSubscriber(quitter), termdash.RedrawInterval(1*time.Second)); err != nil {
			t.logger.Errorf("error running termdash terminal: %s", err)
			// TODO(slok): exit on error.
		}
	}()

	return t.widgets, nil
}

func (t *termDashboard) gridLayout(dashboard model.Dashboard) ([]container.Option, error) {
	builder := grid.New()

	gridElements := []grid.Element{}

	// Create each row.
	rowHeightPerc := (100 / len(dashboard.Rows)) - 1 // All rows same percentage of screen height.
	for _, rowcfg := range dashboard.Rows {

		// TODO(slok): Allow different percentage per widget.
		widgetColPerc := (100 / len(rowcfg.Widgets)) - 1 // All widgets same percentage of screen height.

		// Create widgets per row
		var elements []grid.Element
		for _, widgetcfg := range rowcfg.Widgets {

			var err error
			var widget render.Widget

			// New widget.
			switch {
			case widgetcfg.Gauge != nil:
				widget, err = newGauge(widgetcfg)
				if err != nil {
					t.logger.Errorf("error creating gauge: %s", err)
					continue
				}
			case widgetcfg.Singlestat != nil:
				widget, err = newSinglestat(widgetcfg)
				if err != nil {
					t.logger.Errorf("error creating gauge: %s", err)
					continue
				}
			}

			// Add widget to the tracked widgets so the app can control them.
			t.widgets = append(t.widgets, widget)

			// Create the element and set a size
			element := grid.Widget(widget.(widgetapi.Widget),
				container.Border(linestyle.Light),
				container.BorderTitle(widgetcfg.Title),
			)
			element = grid.ColWidthPerc(widgetColPerc, element)

			// Append to the row.
			elements = append(elements, element)
		}

		rowElement := grid.RowHeightPerc(rowHeightPerc, elements...)
		gridElements = append(gridElements, rowElement)
	}

	// Add rows.
	builder.Add(gridElements...)

	// Get the layout from the grid.
	return builder.Build()
}
