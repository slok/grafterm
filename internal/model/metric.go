package model

import (
	"time"
)

// Metric represents a measured metrics in time.
type Metric struct {
	Value float64
	TS    time.Time
}

// MetricSeries is a group of metrics identified by an ID and a context
// information.
type MetricSeries struct {
	ID      string
	Labels  map[string]string
	Metrics []Metric
}
