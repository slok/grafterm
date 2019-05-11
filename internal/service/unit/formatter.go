package unit

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// Formatter knows how to interact with different values
// to give then a format and a representation.
// They are based on units that can do the conversion to
// other units and represent them in the returning string.
type Formatter func(value float64, decimals int) string

// NewUnitFormatter is a factory that selects the correct formatter
// based on the unit.
// If the unit does not exists it will return an error.
func NewUnitFormatter(unit string) (Formatter, error) {
	unit = strings.ToLower(unit)

	var f Formatter
	switch unit {
	case "", "short":
		f = shortFormatter
	case "none":
		f = noneFormatter
	case "percent":
		f = percentFormatter
	case "ratio":
		f = ratioFormatter
	case "s", "second", "seconds":
		f = secondFormatter
	case "reqps":
		f = newSuffixFormatter(" reqps")
	case "byte", "bytes":
		f = bytesFormatter
	default:
		return nil, fmt.Errorf("%s is not a valid unit", unit)
	}

	return safeFormatter(f), nil
}

// noneFormatter returns the value as it is
// in a float representation with decimal trim.
func noneFormatter(value float64, decimals int) string {
	f := suffixDecimalFormat(decimals, "")
	return fmt.Sprintf(f, value)
}

// percentFormatter returns the value with the
// percent suffix and assumes is a percent value.
// Examples:
//	- 100: 100%
//	- 1029.12: 9876.12%
var percentFormatter = newSuffixFormatter("%")

// ratioFormatter returns the value with the
// percent suffix and assumes is a ratio value
// (0-1).
// Examples:
//	- 1: 100%
//	- 0.412: 41.2%
func ratioFormatter(value float64, decimals int) string {
	return percentFormatter(value*100, decimals)
}

// secondFormatter returns the value with the
// seconds in a single unit pretty format time.
// supports: ns, µs, ms, s, m, h, d.
// Examples:
//	- 1: 1s
//	- 0.1: 100ms
//  - 300: 5m
func secondFormatter(value float64, decimals int) string {
	t := time.Duration(value * float64(time.Second))
	return durationSingleUnitPrettyFormat(t, decimals)
}

// newSuffixFormatter returns a formatter that will apply a
// suffix to the received value.
func newSuffixFormatter(suffix string) Formatter {
	return func(value float64, decimals int) string {
		dFmt := suffixDecimalFormat(decimals, suffix)
		return fmt.Sprintf(dFmt, value)
	}
}

func suffixDecimalFormat(decimals int, suffix string) string {
	suffix = strings.ReplaceAll(suffix, "%", "%%") // Safe `%` character for fmt.
	return fmt.Sprintf("%%.%df%s", decimals, suffix)
}

// durationSingleUnitPrettyFormat returns the pretty format in one single
// unit for a time.Duration, the different returned unit formats
// are: nanoseconds, microseconds, milliseconds, seconds, minutes
// hours, days.
// Implementation obtained from: https://github.com/mum4k/termdash/blob/d34e18ab097be3ec6147767173db12f810f8dbbb/widgets/linechart/value_formatter.go#L31
func durationSingleUnitPrettyFormat(d time.Duration, decimals int) string {
	// Check if the duration is less than 0.
	prefix := ""
	if d < 0 {
		prefix = "-"
		d = time.Duration(math.Abs(d.Seconds()) * float64(time.Second))
	}

	switch {
	// Nanoseconds.
	case d.Nanoseconds() < 1000:
		dFmt := prefix + "%d ns"
		return fmt.Sprintf(dFmt, d.Nanoseconds())
	// Microseconds.
	case d.Seconds()*1000*1000 < 1000:
		dFmt := prefix + suffixDecimalFormat(decimals, " µs")
		return fmt.Sprintf(dFmt, d.Seconds()*1000*1000)
	// Milliseconds.
	case d.Seconds()*1000 < 1000:
		dFmt := prefix + suffixDecimalFormat(decimals, " ms")
		return fmt.Sprintf(dFmt, d.Seconds()*1000)
	// Seconds.
	case d.Seconds() < 60:
		dFmt := prefix + suffixDecimalFormat(decimals, " s")
		return fmt.Sprintf(dFmt, d.Seconds())
	// Minutes.
	case d.Minutes() < 60:
		dFmt := prefix + suffixDecimalFormat(decimals, " m")
		return fmt.Sprintf(dFmt, d.Minutes())
	// Hours.
	case d.Hours() < 24:
		dFmt := prefix + suffixDecimalFormat(decimals, " h")
		return fmt.Sprintf(dFmt, d.Hours())
	// Days.
	default:
		dFmt := prefix + suffixDecimalFormat(decimals, " d")
		return fmt.Sprintf(dFmt, d.Hours()/24)
	}
}

const (
	kibibyte float64 = 1024
	mebibyte         = kibibyte * 1024
	gibibyte         = mebibyte * 1024
	tebibyte         = gibibyte * 1024
	pebibyte         = tebibyte * 1024
	exibyte          = pebibyte * 1024
	zebibyte         = exibyte * 1024
	yobibyte         = zebibyte * 2014
)

// bytesFormatter returns the value  in bytes for a data
// quantity in a pretty format style.
// supports: B, KiB, MiB, GiB, TiB, PiB, EiB, ZiB, YiB.
// Examples:
//	- 35: 35 B
//	- 1024: 1 KiB
var bytesFormatter = newRangedFormatter([]rangeStep{
	{max: kibibyte, base: 1, suffix: " B"},
	{max: mebibyte, base: kibibyte, suffix: " KiB"},
	{max: gibibyte, base: mebibyte, suffix: " MiB"},
	{max: tebibyte, base: gibibyte, suffix: " GiB"},
	{max: pebibyte, base: tebibyte, suffix: " TiB"},
	{max: exibyte, base: pebibyte, suffix: " PiB"},
	{max: zebibyte, base: exibyte, suffix: " EiB"},
	{max: yobibyte, base: zebibyte, suffix: " ZiB"},
	{base: yobibyte, suffix: " YiB"},
})

const (
	shortK     float64 = 1000
	shortMil           = shortK * 1000
	shortBil           = shortMil * 1000
	shortTri           = shortBil * 1000
	shortQuadr         = shortTri * 1000
	shortQuint         = shortQuadr * 1000
	shortSext          = shortQuint * 1000
	shortSept          = shortSext * 1000
)

// shortFormatter returns the value trimming the value on
// high numbers adding a suffix.
// supports: K, Mil, Bil, tri, Quadr, Quint, Sext, Sept.
// Examples:
// 	- 1000 = 1 k
//  - 2000000 = 2 Mil
var shortFormatter = newRangedFormatter([]rangeStep{
	{max: shortK, base: 1, suffix: ""},
	{max: shortMil, base: shortK, suffix: " K"},
	{max: shortBil, base: shortMil, suffix: " Mil"},
	{max: shortTri, base: shortBil, suffix: " Bil"},
	{max: shortQuadr, base: shortTri, suffix: " Tri"},
	{max: shortQuint, base: shortQuadr, suffix: " Quadr"},
	{max: shortSext, base: shortQuint, suffix: " Quint"},
	{max: shortSept, base: shortSext, suffix: " Sext"},
	{base: shortSept, suffix: " Sept"},
})

// safeFormatter wraps a formatter and wraps the received formatter
// with some sanity checks to make it safe.
func safeFormatter(f Formatter) Formatter {
	return func(value float64, decimals int) string {
		if decimals < 0 {
			decimals = 0
		}
		if math.IsNaN(value) {
			return ""
		}
		return f(value, decimals)
	}
}

type rangeStep struct {
	max    float64
	suffix string
	base   float64
}

// newRangeFormatter returns a formatter based on a stepped
// range.
func newRangedFormatter(r []rangeStep) Formatter {
	return func(value float64, decimals int) string {
		// Check if the duration is less than 0.
		prefix := ""
		if value < 0 {
			prefix = "-"
			value = math.Abs(value)
		}

		step := r[0]
		for _, s := range r {
			if value < step.max {
				break
			}
			step = s
		}

		dFmt := prefix + suffixDecimalFormat(decimals, step.suffix)
		return fmt.Sprintf(dFmt, value/step.base)
	}
}
