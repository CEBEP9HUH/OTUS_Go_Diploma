//go:build linux
// +build linux

package linux

import (
	"testing"

	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/services/collector"
	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/services/statistic"
	"github.com/stretchr/testify/require"
)

func TestExtractCPUUsage(t *testing.T) {
	testSuites := map[string]struct {
		input  string
		output statistic.CPUUsage
	}{
		"Simple": {
			input: "5,6 us,  2,8 sy,  0,0 ni, 91,6 id,  0,0 wa,  0,0 hi,  0,0 si,  0,0 st",
			output: statistic.CPUUsage{
				UserMode:   5.6,
				SystemMode: 2.8,
				Idle:       91.6,
			},
		},
		"Reordered": {
			input: "0,0 ni, 91,6 id,  0,0 wa,  5,6 us,  0,0 hi,  0,0 si,  0,0 st,  2,8 sy",
			output: statistic.CPUUsage{
				UserMode:   5.6,
				SystemMode: 2.8,
				Idle:       91.6,
			},
		},
	}
	for name, ts := range testSuites {
		t.Run(name, func(t *testing.T) {
			res, err := extractCPUUsage(ts.input)
			require.NoError(t, err)
			require.EqualValues(t, ts.output, res)
		})
	}
}

func TestExtractCPUUsageErr(t *testing.T) {
	testSuites := map[string]struct {
		input string
		err   error
	}{
		"Empty": {
			input: "",
			err:   collector.ErrCPUInfoNotFound,
		},
		"NoUserMode": {
			input: "2,8 sy,  0,0 ni, 91,6 id",
			err:   collector.ErrUserModeNotFound,
		},
		"BadUserMode": {
			input: "2,8 sy,  bad us, 91,6 id",
			err:   collector.ErrUserModeNotFound,
		},
		"NoSystemMode": {
			input: "2,8 us,  0,0 ni, 91,6 id",
			err:   collector.ErrSystemModeNotFound,
		},
		"BadUSystemMode": {
			input: "2,8 us,  bad sy, 91,6 id",
			err:   collector.ErrSystemModeNotFound,
		},
		"NoIdle": {
			input: "2,8 us,  0,0 ni, 91,6 sy",
			err:   collector.ErrIdleNotFound,
		},
		"BadIdle": {
			input: "2,8 us,  bad id, 91,6 sy",
			err:   collector.ErrIdleNotFound,
		},
	}
	for name, ts := range testSuites {
		t.Run(name, func(t *testing.T) {
			res, err := extractCPUUsage(ts.input)
			require.Nil(t, res)
			require.ErrorIs(t, err, ts.err)
		})
	}
}
