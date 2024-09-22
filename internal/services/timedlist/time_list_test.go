package timedlist

import (
	"sync"
	"sync/atomic"
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
	interval := time.Millisecond
	tl, cancel := MakeTimedList[int](interval)
	defer cancel()
	require.NotNil(t, tl)
	var wg sync.WaitGroup
	var added atomic.Int32

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, v := range vals {
				if tl.Add(v, time.Now()) {
					added.Add(1)
				}
			}
		}()
	}
	wg.Wait()
	assert.Equal(t, added.Load(), int32(tl.Len()))
	<-time.After(interval * 2)
	tl.DeleteExpired()
	assert.Equal(t, 0, tl.Len())
}
