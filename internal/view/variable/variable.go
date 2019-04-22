package variable

import (
	"time"

	"github.com/slok/grafterm/internal/model"
)

type Scope int

const (
	// ScopeDashboard refers to a variable that is used on dashboard load.
	ScopeDashboard Scope = iota
	// ScopeSync refers to a variable that is used on every sync.
	ScopeSync
)

// Variabler represents a variable kind that knows how to get variables.
type Variabler interface {
	Scope() Scope
	// IsRepeatable will return true If the variable is repeatable.
	IsRepeatable() bool
	// GetValue returns the value of the variable.
	// If is a repeatable variable internally it knows
	// how to return the value in one string (e.g `staging|prod|dev`).
	GetValue() string
}

// Repeatable is a variabler that can be repeated.
type Repeatable interface {
	Variabler
	Select(variableID ...string)
	Deselect(variableID ...string)
	GetValues() []string
	GetAllValues() []string
}

// FactoryConfig is the configuration required by the variabler factory.
type FactoryConfig struct {
	TimeRange time.Duration
	Dashboard model.Dashboard
}

// NewVariablers is a factory that knows how to create variablers.
func NewVariablers(cfg FactoryConfig) (map[string]Variabler, error) {
	variablers := map[string]Variabler{}
	for _, v := range cfg.Dashboard.Variables {
		switch {
		case v.Constant != nil:
			variablers[v.Name] = &ConstVariabler{cfg: v}
		case v.Interval != nil:
			variablers[v.Name] = NewIntervalVariabler(cfg.TimeRange, v)
		}
	}

	return variablers, nil
}
