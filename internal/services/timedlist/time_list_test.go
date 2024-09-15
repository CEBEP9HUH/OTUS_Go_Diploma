package timedlist

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreation(t *testing.T) {
	tl, cancel := MakeTimedList[int](time.Second * 3)
	defer cancel()
	require.NotNil(t, tl)
}

func TestAddition(t *testing.T) {
	interval := time.Millisecond
	tl, cancel := MakeTimedList[int](interval)
	defer cancel()
	require.NotNil(t, tl)
	vals := []int{1, 2, 3, 4, 5, 6}

	for _, v := range vals {
		<-time.After(time.Microsecond)
		assert.True(t, tl.Add(v, time.Now()), "Value have to be added")
	}
}

func TestConcurentAddition(t *testing.T) {
	vals := []int{1, 2, 3, 4, 5, 6}
	interval := time.Second
	tl, cancel := MakeTimedList[int](interval)
	defer cancel()
	require.NotNil(t, tl)
	var wg sync.WaitGroup
	for _, v := range []int{5, 7, 11} { // простые числа, чтобы таймауты не пересекались
		wg.Add(1)
		go func(timeout int) {
			defer wg.Done()
			for _, v := range vals {
				<-time.After(time.Duration(v) * time.Microsecond)
				require.True(t, tl.Add(v, time.Now()), "Value have to be added")
			}
		}(v)
	}
	wg.Wait()
}
