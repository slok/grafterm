package termdash

import (
	"context"
	"time"

	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/termbox"
	"github.com/mum4k/termdash/terminal/terminalapi"

	"github.com/slok/meterm/internal/model"
	"github.com/slok/meterm/internal/service/log"
	"github.com/slok/meterm/internal/view/render"
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
	c, err := container.New(
		t.terminal,
		container.Border(linestyle.Light),
		container.BorderTitle("PRESS Q TO QUIT"),
		t.loadDashboard(dashboard),
	)
	if err != nil {
		return nil, err
	}

	quitter := func(k *terminalapi.Keyboard) {
		if k.Key == 'q' || k.Key == 'Q' {
			t.cancel()
		}
	}

	go func() {
		if err := termdash.Run(ctx, t.terminal, c, termdash.KeyboardSubscriber(quitter), termdash.RedrawInterval(1*time.Second)); err != nil {
			t.logger.Errorf("error running termdash terminal: %s", err)
			// TODO(slok): exit on error.
		}
	}()

	return t.widgets, nil
}

func (t *termDashboard) loadDashboard(dashboard model.Dashboard) container.Option {
	return t.loadRow(dashboard.Rows)
}

func (t *termDashboard) loadRow(rows []model.Row) container.Option {
	if len(rows) == 0 {
		return nil
	}

	if len(rows) == 1 {
		return t.loadWidget(rows[0].Widgets)
	}

	// Start splitting.
	spl := len(rows) / 2

	// Top, check if we are on a leaf to add extra opts.
	topSpl := rows[:spl]
	topOpts := []container.Option{}
	if len(topSpl) == 1 {
		topOpts = t.getRowExtraOpts(topSpl[0])
	}
	topOpts = append(topOpts, t.loadRow(topSpl))
	top := container.Top(topOpts...)

	// Bottom, check if we are on a leaf to add extra opts.
	bottomSpl := rows[spl:]
	bottomOpts := []container.Option{}
	if len(bottomSpl) == 1 {
		bottomOpts = t.getRowExtraOpts(bottomSpl[0])
	}
	bottomOpts = append(bottomOpts, t.loadRow(bottomSpl))
	bottom := container.Bottom(bottomOpts...)

	return container.SplitHorizontal(top, bottom)
}

func (t *termDashboard) getRowExtraOpts(row model.Row) []container.Option {
	extraOpts := []container.Option{}
	if row.Title != "" {
		extraOpts = append(extraOpts, container.BorderTitle(row.Title))
	}

	if row.Border {
		extraOpts = append(extraOpts, container.Border(linestyle.Light))
	}
	return extraOpts
}

func (t *termDashboard) loadWidget(ws []model.Widget) container.Option {
	if len(ws) == 0 {
		return nil
	}

	if len(ws) == 1 {
		// New widget.
		// TODO(slok): Check other widgets.
		cfg := ws[0]
		w, err := newGauge(cfg)
		if err != nil {
			t.logger.Errorf("error creating gauge: %s", err)
		}
		t.widgets = append(t.widgets, w)

		return container.PlaceWidget(w)
	}

	// Start splitting.
	spl := len(ws) / 2
	return container.SplitVertical(
		container.Left(
			t.loadWidget(ws[:spl]),
		),
		container.Right(
			t.loadWidget(ws[spl:]),
		),
	)
}
