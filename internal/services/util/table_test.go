package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetHeadsIds(t *testing.T) {
	testSuites := map[string]struct {
		header []string
		heads  []string
		out    map[string]int
	}{
		"TestSimple": {
			header: []string{"aaa", "bbb", "ccc"},
			heads:  []string{"aaa", "ccc"},
			out:    map[string]int{"aaa": 0, "ccc": 2},
		},
		"TestSpaces": {
			header: []string{"aaa", "bbb", "c c"},
			heads:  []string{"aaa", "c c"},
			out:    map[string]int{"aaa": 0, "c c": 2},
		},
		"TestCases": {
			header: []string{"aaa", "bbb", "c C"},
			heads:  []string{"aaa", "c C"},
			out:    map[string]int{"aaa": 0, "c C": 2},
		},
		"TestEmpty": {
			header: []string{"aaa", "bbb", "c C"},
			heads:  []string{},
			out:    map[string]int{},
		},
		"TestMixed": {
			header: []string{"aaa", "bbb", "c C"},
			heads:  []string{"c C", "aaa"},
			out:    map[string]int{"aaa": 0, "c C": 2},
		},
	}

	for name, ts := range testSuites {
		t.Run(name, func(t *testing.T) {
			out, err := GetHeadsIDs(ts.header, ts.heads)
			assert.NoError(t, err, "Error is unexpected")

			for k, expectedV := range ts.out {
				v, ok := out[k]
				require.True(t, ok, "Result have to have value with key %q", k)
				assert.EqualValues(t, expectedV, v)
			}
		})
	}
}

func TestGetHeadsIdsErr(t *testing.T) {
	testSuites := map[string]struct {
		header []string
		heads  []string
	}{
		"TestSimple": {
			header: []string{"aaa", "bbb", "ccc"},
			heads:  []string{"ddd", "eee"},
		},
		"TestPartialMatch": {
			header: []string{"aaa", "bbb", "ccc"},
			heads:  []string{"ddd", "ccc"},
		},
		"TestTooManyHeads": {
			header: []string{"aaa", "bbb", "ccc"},
			heads:  []string{"aaa", "bbb", "ccc", "ddd"},
		},
		"TestReapetedRequestedHead": {
			header: []string{"aaa", "bbb", "c C"},
			heads:  []string{"c C", "aaa", "c C"},
		},
		"TestReapetedInputHead": {
			header: []string{"aaa", "bbb", "c C", "aaa"},
			heads:  []string{"c C", "aaa"},
		},
	}

	for name, ts := range testSuites {
		t.Run(name, func(t *testing.T) {
			out, err := GetHeadsIDs(ts.header, ts.heads)
			assert.ErrorIs(t, err, ErrHeadersNotFound, "Error is expected")
			assert.Len(t, out, 0)
		})
	}
}
