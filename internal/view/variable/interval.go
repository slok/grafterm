package variable

import (
	"time"

	"github.com/slok/meterm/internal/model"
	"github.com/slok/meterm/internal/service/unit"
)

type intervalVariabler struct {
	intervalStr string
	cfg         model.Variable
}

// NewIntervalVariabler returns a new variabler that knows how to set
// variables based on the interval, at this moment it only returns
// autoinverval so is not repeatable.
// TODO(slok): make repeatable and allow selecting multiple intervals.
func NewIntervalVariabler(timeRange time.Duration, cfg model.Variable) Variabler {
	// Set default auto interval if not 0.
	steps := 50
	if cfg.Interval.Steps != 0 {
		steps = cfg.Interval.Steps
	}
	dur := unit.NearestDurationFromSteps(timeRange, steps)
	durStr := unit.DurationToSimpleString(dur)

	return &intervalVariabler{
		cfg:         cfg,
		intervalStr: durStr,
	}
}

func (i intervalVariabler) Scope() Scope {
	return ScopeDashboard
}

func (i intervalVariabler) IsRepeatable() bool {
	return false
}

func (i intervalVariabler) GetValue() string {
	return i.intervalStr
}
