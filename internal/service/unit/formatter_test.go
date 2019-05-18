package unit_test

import (
	"math"
	"testing"

	"github.com/slok/grafterm/internal/service/unit"
	"github.com/stretchr/testify/assert"
)

func TestUnitFormatter(t *testing.T) {
	tests := map[string]struct {
		unit     string
		value    float64
		decimals int
		expStr   string
		expErr   bool
	}{
		// Special cases.
		"Not valid.":        {unit: "unknown", expErr: true},
		"Invalid decimals.": {unit: "none", value: 1234.12345, decimals: -10, expStr: "1234"},
		"NaN.":              {unit: "", value: math.NaN(), expStr: ""},

		// No format.
		"Default with no decimals.": {unit: "", value: 1234.12345, decimals: 0, expStr: "1 K"},
		"Default with decimals.":    {unit: "", value: 1234.12345, decimals: 2, expStr: "1.23 K"},

		// None.
		"None with no decimals.": {unit: "none", value: 1234.12345, decimals: 0, expStr: "1234"},
		"None with decimals.":    {unit: "none", value: 1234.12345, decimals: 1, expStr: "1234.1"},

		// Percent.
		"Percent with no decimals.": {unit: "percent", value: 1234.12345, decimals: 0, expStr: "1234%"},
		"Percent with 2 decimals.":  {unit: "percent", value: 1234.12345, decimals: 2, expStr: "1234.12%"},

		// Ratio.
		"Ratio with no decimals.": {unit: "ratio", value: 1.12345, decimals: 0, expStr: "112%"},
		"Ratio with 2 decimals.":  {unit: "ratio", value: 0.12345, decimals: 2, expStr: "12.35%"},

		// SecondsDuration.
		"Seconds with 0.":                      {unit: "seconds", value: 0, decimals: 0, expStr: "0 ns"},
		"Seconds with minus and no decimals.":  {unit: "seconds", value: -1500, decimals: 0, expStr: "-25 m"},
		"Seconds to ns with no decimals.":      {unit: "seconds", value: 1.23e-7, decimals: 0, expStr: "123 ns"},
		"Seconds to µs with no decimals.":      {unit: "seconds", value: 1.23e-4, decimals: 0, expStr: "123 µs"},
		"Seconds to ms with no decimals.":      {unit: "seconds", value: 0.123, decimals: 0, expStr: "123 ms"},
		"Seconds to seconds with no decimals.": {unit: "seconds", value: 35.1234, decimals: 0, expStr: "35 s"},
		"Seconds to m with no decimals.":       {unit: "seconds", value: 60.1, decimals: 0, expStr: "1 m"},
		"Seconds to h with no decimals.":       {unit: "seconds", value: 2 * 60 * 60, decimals: 0, expStr: "2 h"},
		"Seconds to d with no decimals.":       {unit: "seconds", value: 5 * 24 * 60 * 60, decimals: 0, expStr: "5 d"},
		"Seconds with minus and decimals.":     {unit: "seconds", value: -1535, decimals: 2, expStr: "-25.58 m"},
		"Seconds to ns with decimals.":         {unit: "seconds", value: 1.23e-7, decimals: 3, expStr: "123 ns"},
		"Seconds to µs with decimals.":         {unit: "seconds", value: 1.2323e-4, decimals: 5, expStr: "123.23000 µs"},
		"Seconds to ms with decimals.":         {unit: "seconds", value: 0.123098765, decimals: 7, expStr: "123.0987650 ms"},
		"Seconds to seconds with decimals.":    {unit: "seconds", value: 35.1234, decimals: 2, expStr: "35.12 s"},
		"Seconds to m with decimals.":          {unit: "seconds", value: 64, decimals: 1, expStr: "1.1 m"},
		"Seconds to h with decimals.":          {unit: "seconds", value: 2*60*60 + 47, decimals: 4, expStr: "2.0131 h"},
		"Seconds to d with decimals.":          {unit: "seconds", value: 5*24*60*60 + 92012, decimals: 6, expStr: "6.064954 d"},

		// MillisecondsDuration.
		"Milliseconds with 0.":                      {unit: "milliseconds", value: 0, decimals: 0, expStr: "0 ns"},
		"Milliseconds with minus and no decimals.":  {unit: "milliseconds", value: -1500, decimals: 0, expStr: "-2 s"},
		"Milliseconds to minutes with decimals.":    {unit: "milliseconds", value: 150*1000 + 47, decimals: 4, expStr: "2.5008 m"},

		// Request/s.
		"Reqps with no decimals.": {unit: "reqps", value: 1.12345, decimals: 0, expStr: "1 reqps"},
		"Reqps with 2 decimals.":  {unit: "reqps", value: 0.12345, decimals: 2, expStr: "0.12 reqps"},

		// Bytes.
		"Bytes with 0.":                        {unit: "bytes", value: 0, decimals: 0, expStr: "0 B"},
		"Bytes to bytes with no decimals.":     {unit: "bytes", value: -107, decimals: 0, expStr: "-107 B"},
		"Bytes to kibibytes with no decimals.": {unit: "bytes", value: 2050, decimals: 0, expStr: "2 KiB"},
		"Bytes to kibibytes with decimals.":    {unit: "bytes", value: 2081, decimals: 2, expStr: "2.03 KiB"},
		"Bytes to mebibytes with no decimals.": {unit: "bytes", value: 1.405e+8, decimals: 0, expStr: "134 MiB"},
		"Bytes to mebibytes with decimals.":    {unit: "bytes", value: 1.405e+8, decimals: 3, expStr: "133.991 MiB"},
		"Bytes to gibibytes with no decimals.": {unit: "bytes", value: 6.034e+11, decimals: 0, expStr: "562 GiB"},
		"Bytes to gibibytes with decimals.":    {unit: "bytes", value: 6.034e+11, decimals: 6, expStr: "561.960042 GiB"},
		"Bytes to tebibytes with no decimals.": {unit: "bytes", value: 4.508e+13, decimals: 0, expStr: "41 TiB"},
		"Bytes to tebibytes with decimals.":    {unit: "bytes", value: 4.812e+13, decimals: 3, expStr: "43.765 TiB"},
		"Bytes to pebibytes with no decimals.": {unit: "bytes", value: 1.914e+16, decimals: 0, expStr: "17 PiB"},
		"Bytes to pebibytes with decimals.":    {unit: "bytes", value: 1.934e+16, decimals: 3, expStr: "17.177 PiB"},
		"Bytes to exibytes with no decimals.":  {unit: "bytes", value: 4.812e+19, decimals: 0, expStr: "42 EiB"},
		"Bytes to exibytes with decimals.":     {unit: "bytes", value: 4.812e+19, decimals: 2, expStr: "41.74 EiB"},
		"Bytes to zebibytes with no decimals.": {unit: "bytes", value: 15.72e+21, decimals: 0, expStr: "13 ZiB"},
		"Bytes to zebibytes with decimals.":    {unit: "bytes", value: 15.72e+21, decimals: 1, expStr: "13.3 ZiB"},
		"Bytes to yebibytes with no decimals.": {unit: "bytes", value: 862.12e+23, decimals: 0, expStr: "36 YiB"},
		"Bytes to yebibytes with decimals.":    {unit: "bytes", value: 862.12e+23, decimals: 9, expStr: "36.258393470 YiB"},

		// Short.
		"Short with 0.":                    {unit: "short", value: 0, decimals: 0, expStr: "0"},
		"Short to short with no decimals.": {unit: "short", value: -107, decimals: 0, expStr: "-107"},
		"Short to short with decimals.":    {unit: "short", value: 234.12345, decimals: 4, expStr: "234.1234"},
		"Short to K with no decimals.":     {unit: "short", value: 2050, decimals: 0, expStr: "2 K"},
		"Short to K with decimals.":        {unit: "short", value: 2050, decimals: 2, expStr: "2.05 K"},
		"Short to Mil with no decimals.":   {unit: "short", value: 1.56e+6, decimals: 0, expStr: "2 Mil"},
		"Short to Mil with decimals.":      {unit: "short", value: 1.56e+6, decimals: 3, expStr: "1.560 Mil"},
		"Short to Bil with no decimals.":   {unit: "short", value: 5.783e+9, decimals: 0, expStr: "6 Bil"},
		"Short to Bil with decimals.":      {unit: "short", value: 5.783e+9, decimals: 6, expStr: "5.783000 Bil"},
		"Short to Tri with no decimals.":   {unit: "short", value: 82.321e+12, decimals: 0, expStr: "82 Tri"},
		"Short to Tri with decimals.":      {unit: "short", value: 82.321e+12, decimals: 1, expStr: "82.3 Tri"},
		"Short to Quadr with no decimals.": {unit: "short", value: 9.6132e+15, decimals: 0, expStr: "10 Quadr"},
		"Short to Quadr with decimals.":    {unit: "short", value: 9.6132e+15, decimals: 4, expStr: "9.6132 Quadr"},
		"Short to Quint with no decimals.": {unit: "short", value: 31.99e+18, decimals: 0, expStr: "32 Quint"},
		"Short to Quint with decimals.":    {unit: "short", value: 31.99e+18, decimals: 2, expStr: "31.99 Quint"},
		"Short to Sext with no decimals.":  {unit: "short", value: 17.581e+21, decimals: 0, expStr: "18 Sext"},
		"Short to Sext with decimals.":     {unit: "short", value: 17.581e+21, decimals: 5, expStr: "17.58100 Sext"},
		"Short to Sept with no decimals.":  {unit: "short", value: 2.812e+25, decimals: 0, expStr: "28 Sept"},
		"Short to Sept with decimals.":     {unit: "short", value: 2.812e+25, decimals: 5, expStr: "28.12000 Sept"},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			f, err := unit.NewUnitFormatter(test.unit)
			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				gotStr := f(test.value, test.decimals)
				assert.Equal(test.expStr, gotStr)
			}
		})
	}
}
