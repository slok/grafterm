package model

import (
	"time"
)

// Metric represents a measured metrics in time.
type Metric struct {
	Value float64
	TS    time.Time
}

// Series is a group of metrics identified by an ID and a context
// information.
type Series struct {
	ID      string
	Labels  map[string]string
	Metrics []Metric
}
