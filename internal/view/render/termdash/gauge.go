package termdash

import (
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgets/donut"

	"github.com/slok/meterm/internal/model"
)

// gauge satisfies render.GaugeWidget interface.
type gauge struct {
	cfg model.Widget
	*donut.Donut
}

func newGauge(cfg model.Widget) (*gauge, error) {
	// Create the widget.
	donut, err := donut.New(donut.CellOpts(cell.FgColor(cell.ColorWhite)))
	if err != nil {
		return nil, err
	}

	return &gauge{
		Donut: donut,
		cfg:   cfg,
	}, nil
}

func (g *gauge) GetWidgetCfg() model.Widget {
	return g.cfg
}

func (g *gauge) Sync(isPercent bool, value float64) error {
	var err error
	if isPercent {
		err = g.Donut.Percent(int(value))
	} else {
		max := float64(g.cfg.Gauge.Max)
		if max < value {
			max = value
		}
		err = g.Donut.Absolute(int(value), int(max))
	}

	if err != nil {
		return err
	}

	return nil
}

func (g *gauge) SetColor(hexColor string) error {
	color, err := colorHexToTermdash(hexColor)
	if err != nil {
		return err
	}

	// Create a new widget with the current color.
	d, err := donut.New(donut.CellOpts(cell.FgColor(color)))
	if err != nil {
		return err
	}

	// Replace the instance.
	g.Donut = d

	return nil
}
