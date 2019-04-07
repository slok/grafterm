package termdash

import (
	"fmt"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/mum4k/termdash/cell"
)

func colorHexToTermdash(color string) (cell.Color, error) {
	c, err := colorful.Hex(color)
	if err != nil {
		return 0, fmt.Errorf("error getting color: %s", err)
	}

	cr, cg, cb := c.RGB255()
	return cell.ColorRGB24(int(cr), int(cg), int(cb)), nil
}
