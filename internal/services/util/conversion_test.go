package util

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	delta = 1e-9
)

var testSuites = map[string]struct {
	in  string
	out float32
}{
	"TestSimpleConversion1":  {"1.2", float32(1.2)},
	"TestSimpleConversion2":  {"1.2345", float32(1.2345)},
	"TestDashConversion":     {"-", float32(0)},
	"TestCommaConversion1":   {"1,2", float32(1.2)},
	"TestCommaConversion2":   {"1,2345", float32(1.2345)},
	"TestNegativeConversion": {"-1,2345", float32(-1.2345)},
	"TestZeroConversion":     {"0", float32(0)},
	"TestShortFormat1":       {".23", float32(0.23)},
	"TestShortFormat2":       {",23", float32(0.23)},
	"TestShortFormat3":       {"23.", float32(23)},
	"TestShortFormat4":       {"23,", float32(23)},
}

var badTestSuites = map[string]string{
	"TestTextConversion":            "asd",
	"TestMultipleDashesConversion1": "----",
	"TestMultipleDashesConversion2": "--1.2",
}

func TestGetFloat(t *testing.T) {
	for name, ts := range testSuites {
		t.Run(name, func(t *testing.T) {
			out, err := GetFloat(ts.in)
			assert.NoError(t, err, "Error is unexpected")
			assert.InDelta(t, ts.out, out, delta, "String to float conversion failed")
		})
	}
}

func TestGetFloatErr(t *testing.T) {
	for name, ts := range badTestSuites {
		t.Run(name, func(t *testing.T) {
			_, err := GetFloat(ts)
			assert.Error(t, err, "Error is expected")
		})
	}
}

func TestGetTrimmedFloat(t *testing.T) {
	for name, ts := range testSuites {
		t.Run(name, func(t *testing.T) {
			out, err := GetTrimmedFloat(ts.in, "", "")
			assert.NoError(t, err, "Error is unexpected")
			assert.InDelta(t, ts.out, out, delta, "String to float conversion failed")
		})
	}
}

func TestGetTrimmedFloat_PrefixAndSuffix(t *testing.T) {
	const (
		prefix = "prefix"
		suffix = "suffix"
	)
	localTestSuites := make(map[string]struct {
		in  string
		out float32
	}, len(testSuites))
	for k, v := range testSuites {
		localTestSuites[k] = struct {
			in  string
			out float32
		}{
			in:  fmt.Sprintf("%s%s%s", prefix, v.in, suffix),
			out: v.out,
		}
	}
	for name, ts := range localTestSuites {
		t.Run(name, func(t *testing.T) {
			out, err := GetTrimmedFloat(ts.in, prefix, suffix)
			assert.NoError(t, err, "Error is unexpected")
			assert.InDelta(t, ts.out, out, delta, "String to float conversion failed")
		})
	}
}

func TestGetTrimmedFloatErr(t *testing.T) {
	for name, ts := range badTestSuites {
		t.Run(name, func(t *testing.T) {
			_, err := GetFloat(ts)
			assert.Error(t, err, "Error is expected")
		})
	}
}

func TestGetTrimmedFloatErr_PrefixAndSuffix(t *testing.T) {
	const (
		prefix = "prefix"
		suffix = "suffix"
	)
	localTestSuites := make(map[string]string, len(badTestSuites))
	for k, v := range badTestSuites {
		localTestSuites[k] = fmt.Sprintf("%s%s%s", prefix, v, suffix)
	}
	for name, ts := range localTestSuites {
		t.Run(name, func(t *testing.T) {
			_, err := GetTrimmedFloat(ts, prefix, suffix)
			assert.Error(t, err, "Error is expected")
		})
	}
}
