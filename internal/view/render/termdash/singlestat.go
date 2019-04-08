package termdash

import (
	"fmt"

	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgets/segmentdisplay"

	"github.com/slok/meterm/internal/model"
)

const (
	defFormat = "%.2f"
)

// singlestat satisfies render.SinglestatWidget interface.
type singlestat struct {
	cfg     model.Widget
	color   cell.Color
	textFmt string
	*segmentdisplay.SegmentDisplay
}

func newSinglestat(cfg model.Widget) (*singlestat, error) {
	// Create the widget.
	sd, err := segmentdisplay.New()
	if err != nil {
		return nil, err
	}

	// Create the segment display format.
	textFmt := defFormat
	if cfg.Singlestat.TextFormat != "" {
		textFmt = cfg.Singlestat.TextFormat
	}

	return &singlestat{
		SegmentDisplay: sd,
		textFmt:        textFmt,
		color:          cell.ColorWhite,
		cfg:            cfg,
	}, nil
}

func (s *singlestat) GetWidgetCfg() model.Widget {
	return s.cfg
}

func (s *singlestat) Sync(value float64) error {
	chunks := []*segmentdisplay.TextChunk{
		segmentdisplay.NewChunk(
			fmt.Sprintf(s.textFmt, value),
			segmentdisplay.WriteCellOpts(cell.FgColor(s.color))),
	}
	err := s.Write(chunks)
	if err != nil {
		return err
	}
	return nil
}

func (s *singlestat) SetColor(hexColor string) error {
	color, err := colorHexToTermdash(hexColor)
	if err != nil {
		return err
	}
	s.color = color
	return nil
}
