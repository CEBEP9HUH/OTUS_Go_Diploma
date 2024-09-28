//go:build linux
// +build linux

package linux

import (
	"testing"

	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/services/collector"
	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/services/statistic"
	"github.com/stretchr/testify/require"
)

func TestExtractLoadAvgUsage(t *testing.T) {
	testSuites := map[string]struct {
		input  string
		output statistic.LoadAvg
	}{
		"Simple": {
			input: "0,95, 0,85, 0,75",
			output: statistic.LoadAvg{
				Min1:  0.95,
				Min5:  0.85,
				Min15: 0.75,
			},
		},
	}
	for name, ts := range testSuites {
		t.Run(name, func(t *testing.T) {
			res, err := extractloadAvgUsage(ts.input)
			require.NoError(t, err)
			require.EqualValues(t, ts.output, res)
		})
	}
}

func TestExtractLoadAvgUsageErr(t *testing.T) {
	testSuites := map[string]struct {
		input string
		err   error
	}{
		"Empty": {
			input: "",
			err:   collector.ErrLoadAvgInfoNotFound,
		},
		"Not Full": {
			input: "22 33",
			err:   collector.ErrLoadAvgInfoNotFound,
		},
		"Not Found 1": {
			input: "bad, 0,23, 0,1",
			err:   collector.ErrAvgMin1NotFound,
		},
		"Not Found 5": {
			input: "0,1, bad, 0,23",
			err:   collector.ErrAvgMin5NotFound,
		},
		"Not Found 15": {
			input: "0,1, 0,23, bad",
			err:   collector.ErrAvgMin15NotFound,
		},
		"Bad fromat #1": {
			input: "0,1,0,23,0,24",
			err:   collector.ErrLoadAvgInfoNotFound,
		},
		"Bad fromat #2": {
			input: "0,95, 0,85, 0,75, 0,123",
			err:   collector.ErrLoadAvgInfoNotFound,
		},
	}
	for name, ts := range testSuites {
		t.Run(name, func(t *testing.T) {
			res, err := extractloadAvgUsage(ts.input)
			require.Nil(t, res)
			require.ErrorIs(t, err, ts.err)
		})
	}
}
