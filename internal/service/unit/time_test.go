package unit_test

import (
	"testing"
	"time"

	"github.com/slok/meterm/internal/service/unit"
	"github.com/stretchr/testify/assert"
)

func TestNearestDurationFromSteps(t *testing.T) {
	tests := []struct {
		name      string
		timeRange time.Duration
		steps     int
		expDur    time.Duration
	}{
		{
			name:      "A low range with lots of steps should return the min interval",
			timeRange: 2 * time.Hour,
			steps:     1000,
			expDur:    30 * time.Second,
		},
		{
			name:      "Greater that the first available interval but not greater than the second one.",
			timeRange: 31 * time.Second,
			steps:     1,
			expDur:    30 * time.Second,
		},
		{
			name:      "A high range with few steps should return the max interval",
			timeRange: 5000 * time.Hour,
			steps:     2,
			expDur:    720 * time.Hour,
		},
		{
			name:      "calculates the nearest (12h, 50 steps).",
			timeRange: 12 * time.Hour,
			steps:     50,
			expDur:    10 * time.Minute,
		},
		{
			name:      "calculates the nearest (2h, 50 steps).",
			timeRange: 2 * time.Hour,
			steps:     50,
			expDur:    2 * time.Minute,
		},
		{
			name:      "calculates the nearest (30m, 50 steps).",
			timeRange: 30 * time.Minute,
			steps:     50,
			expDur:    30 * time.Second,
		},
		{
			name:      "calculates the nearest (6h, 30 steps).",
			timeRange: 6 * time.Hour,
			steps:     30,
			expDur:    10 * time.Minute,
		},
		{
			name:      "calculates the nearest (3d, 50 steps).",
			timeRange: 72 * time.Hour,
			steps:     50,
			expDur:    1 * time.Hour,
		},
		{
			name:      "calculates the nearest (7d, 50 steps).",
			timeRange: 168 * time.Hour,
			steps:     50,
			expDur:    3 * time.Hour,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotDur := unit.NearestDurationFromSteps(test.timeRange, test.steps)
			assert.Equal(t, test.expDur, gotDur)
		})
	}
}

func TestDurationToSimpleString(t *testing.T) {
	tests := []struct {
		name      string
		timeRange time.Duration
		exp       string
	}{
		{
			name:      "Seconds.",
			timeRange: 39 * time.Second,
			exp:       "39s",
		},
		{
			name:      "Minutes.",
			timeRange: 39 * time.Minute,
			exp:       "39m",
		},
		{
			name:      "Minutes in seconds.",
			timeRange: 61 * time.Second,
			exp:       "61s",
		},
		{
			name:      "Hours.",
			timeRange: 39 * time.Hour,
			exp:       "39h",
		},
		{
			name:      "Hours in minutes.",
			timeRange: 75 * time.Minute,
			exp:       "75m",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := unit.DurationToSimpleString(test.timeRange)
			assert.Equal(t, test.exp, got)
		})
	}
}

func TestSteppedTimeRangeStringFormat(t *testing.T) {
	tests := []struct {
		name      string
		timeRange time.Duration
		steps     int
		exp       string
	}{
		{
			name:      "Month based ranges should have the day and month.",
			timeRange: 38 * 24 * time.Hour,
			exp:       "01/02",
		},
		{
			name:      "More than a half of a month ranges should have the day and month.",
			timeRange: 16 * 24 * time.Hour,
			exp:       "01/02",
		},
		{
			name:      "More than one day ranges should have day and time.",
			timeRange: 72 * time.Hour,
			exp:       "01/02 15:04",
		},
		{
			name:      "Less than a minute ranges should have seconds.",
			timeRange: 48 * time.Second,
			exp:       "15:04:05",
		},
		{
			name:      "less than one day based ranges with few steps don't have seconds.",
			timeRange: 2 * time.Hour,
			steps:     50,
			exp:       "15:04",
		},
		{
			name:      "less than one day based ranges with few steps don't have seconds.",
			timeRange: 15 * time.Minute,
			steps:     200,
			exp:       "15:04:05",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := unit.TimeRangeTimeStringFormat(test.timeRange, test.steps)
			assert.Equal(t, test.exp, got)
		})
	}
}
