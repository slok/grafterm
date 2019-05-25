package widget

import (
	"fmt"
	"sync"

	"github.com/slok/grafterm/internal/model"
)

// Default colors got from Grafana.
// Check: https://github.com/grafana/grafana/blob/406ef962fc113091bb7229c8f3f0090d63c8392e/packages/grafana-ui/src/utils/colors.ts#L16
var defColors = []string{
	"#7EB26D", // 0: pale green
	"#EAB839", // 1: mustard
	"#6ED0E0", // 2: light blue
	"#EF843C", // 3: orange
	"#E24D42", // 4: red
	"#1F78C1", // 5: ocean
	"#BA43A9", // 6: purple
	"#705DA0", // 7: violet
	"#508642", // 8: dark green
	"#CCA300", // 9: dark sand
	"#447EBC",
	"#C15C17",
	"#890F02",
	"#0A437C",
	"#6D1F62",
	"#584477",
	"#B7DBAB",
	"#F4D598",
	"#70DBED",
	"#F9BA8F",
	"#F29191",
	"#82B5D8",
	"#E5A8E2",
	"#AEA2E0",
	"#629E51",
	"#E5AC0E",
	"#64B0C8",
	"#E0752D",
	"#BF1B00",
	"#0A50A1",
	"#962D82",
	"#614D93",
	"#9AC48A",
	"#F2C96D",
	"#65C5DB",
	"#F9934E",
	"#EA6460",
	"#5195CE",
	"#D683CE",
	"#806EB7",
	"#3F6833",
	"#967302",
	"#2F575E",
	"#99440A",
	"#58140C",
	"#052B51",
	"#511749",
	"#3F2B5B",
	"#E0F9D7",
	"#FCEACA",
	"#CFFAFF",
	"#F9E2D2",
	"#FCE2DE",
	"#BADFF4",
	"#F9D9F9",
	"#DEDAF7",
}

type syncingFlag struct {
	syncing bool
	mu      sync.Mutex
}

// Set will return true if it has changed the value and false if already
// was on that state, this way the setter knows if other part of the app has
// changed in the interval it was calling set.
func (s *syncingFlag) Set(v bool) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.syncing == v {
		return false
	}

	s.syncing = v
	return true
}

func (s *syncingFlag) Get() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.syncing
}

// widgetColorManager manages the color selection for widgets.
// it knows to get default color, based on series legend...
// The color selector tracks the number of default colors returned
// so it doesn't repeat default colors.
type widgetColorManager struct {
	count int
}

// GetColorFromSeriesLegend will return the configured color for the matching regex with the series
// legend, if there is no match then it will return a default color.
func (w *widgetColorManager) GetColorFromSeriesLegend(cfg model.GraphWidgetSource, legend string) string {
	so, ok := seriesOverride(cfg.Visualization.SeriesOverride, legend)
	if ok && so.Color != "" {
		return so.Color
	}

	// No match, get the next default color,
	return w.GetDefaultColor()
}

// GetColorFromThresholds gets the correct color based on a ordered list of thresholds and a value.
func (w widgetColorManager) GetColorFromThresholds(thresholds []model.Threshold, value float64) (hexColor string, err error) {
	if len(thresholds) == 0 {
		return "", fmt.Errorf("the number of thresholds can't be 0")
	}

	// Search the correct color.
	threshold := thresholds[0]
	for _, t := range thresholds[1:] {
		if value >= t.StartValue {
			threshold = t
		}
	}

	return threshold.Color, nil
}

// GetDefaultColor returns a default color, for each returned default color the manager
// will track how many default colors have been returned so it doesn't repeat until all
// the default color list has been used and it starts again from the first default color.
func (w *widgetColorManager) GetDefaultColor() string {
	color := defColors[w.count]
	w.count++
	if w.count >= len(defColors) {
		w.count = 0
	}

	return color
}

// seriesOverride returns the series override based on the series legend
// if it finds one, if not then it will return false in the ok return
// argument.
func seriesOverride(seriesOverride []model.SeriesOverride, legend string) (so model.SeriesOverride, ok bool) {
	for _, so := range seriesOverride {
		if so.CompiledRegex != nil && so.CompiledRegex.MatchString(legend) {
			return so, true
		}
	}

	return model.SeriesOverride{}, false
}
