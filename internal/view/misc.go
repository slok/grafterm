package view

import (
	"fmt"
	"sync"

	"github.com/slok/meterm/internal/model"
)

// getThresholdColor gets the correct color based on a ordered list of thresholds
// and a value.
func getThresholdColor(thresholds []model.Threshold, value float64) (hexColor string, err error) {
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
