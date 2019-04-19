package unit

import (
	"fmt"
	"time"
)

var (
	defIntervals = []time.Duration{
		30 * time.Second,
		1 * time.Minute,
		2 * time.Minute,
		10 * time.Minute,
		30 * time.Minute,
		1 * time.Hour,
		3 * time.Hour,
		6 * time.Hour,
		12 * time.Hour,
		24 * time.Hour,
		168 * time.Hour, // 7d
		336 * time.Hour, // 14d
		720 * time.Hour, // 30d
	}
)

// NearestDurationFromSteps returns the nearest interval based on the
// steps in a time range.
func NearestDurationFromSteps(timeRange time.Duration, steps int) time.Duration {
	rawInterval := timeRange / time.Duration(steps)

	switch {
	case rawInterval <= defIntervals[0]:
		return defIntervals[0]
	case rawInterval >= defIntervals[len(defIntervals)-1]:
		return defIntervals[len(defIntervals)-1]
	}

	return getNearestDuration(defIntervals, rawInterval)
}

func getNearestDuration(intervals []time.Duration, timeRange time.Duration) time.Duration {
	var bottom, top time.Duration

	// Get the top and bottom limits in the range.
	for _, limit := range intervals {
		if limit <= timeRange {
			bottom = limit
			continue
		}

		top = limit
		break
	}

	// Get distance from both and return the shortest one.
	bottomDiff := timeRange - bottom
	topDiff := top - timeRange
	if bottomDiff < topDiff {
		return bottom
	}
	return top
}

// DurationToSimpleString will get a duration interval and get the string
// with a simple format (e.g 14m instead of 14m0s).
func DurationToSimpleString(dur time.Duration) string {
	res := ""
	switch {
	case int(dur.Minutes()) < 1,
		int(dur.Seconds())%60 != 0:
		res = fmt.Sprintf("%.0fs", dur.Seconds())
	case int(dur.Hours()) < 1,
		int(dur.Minutes())%60 != 0:
		res = fmt.Sprintf("%.0fm", dur.Minutes())
	default:
		res = fmt.Sprintf("%.0fh", dur.Hours())
	}

	return res
}
