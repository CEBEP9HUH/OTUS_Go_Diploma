package timedlist

import (
	"sync"
	"sync/atomic"
	"time"
)

type node[T any] struct {
	timestamp time.Time
	data      T
	next      *node[T]
	prev      *node[T]
}

type timedList[T any] struct {
	storeDuration time.Duration

	rwm sync.RWMutex

	len    atomic.Int32
	first  *node[T]
	last   *node[T]
	cancel <-chan struct{}
}

type TimedList[T any] interface {
	Len() int
	StoreDuration() time.Duration
	StoredInterval() time.Duration
	GetAfter(t time.Time) ([]T, time.Time)
	GetLast(t time.Duration) ([]T, time.Time)

	SetStoreDuration(storeDuration time.Duration)
	Add(T, time.Time) bool
	DeleteExpired()
}

func MakeTimedList[T any](storeDuration time.Duration) (TimedList[T], func()) {
	cancel := make(chan struct{})
	ret := timedList[T]{
		storeDuration: storeDuration,
		cancel:        cancel,
	}
	go func(tl TimedList[T]) {
		for {
			select {
			case <-time.After(storeDuration):
				tl.DeleteExpired()
			case <-cancel:
				return
			}
		}
	}(&ret)
	return &ret, func() { close(cancel) }
}

func (tl *timedList[T]) Len() int {
	return int(tl.len.Load())
}

func (tl *timedList[T]) SetStoreDuration(storeDuration time.Duration) {
	tl.rwm.Lock()
	defer tl.rwm.Unlock()
	tl.storeDuration = storeDuration
}

func (tl *timedList[T]) StoreDuration() time.Duration {
	tl.rwm.RLock()
	defer tl.rwm.RUnlock()
	return tl.storeDuration
}

func (tl *timedList[T]) Add(data T, t time.Time) bool {
	tl.rwm.Lock()
	defer tl.rwm.Unlock()
	n := node[T]{
		timestamp: t,
		data:      data,
	}
	if tl.last == nil {
		tl.last = &n
		tl.first = tl.last
		tl.len.Store(1)
		return true
	}

	defer tl.len.Add(1)
	pos := tl.last
	for {
		if pos.timestamp.Equal(n.timestamp) {
			return false
		}
		if pos.timestamp.Before(n.timestamp) {
			break
		}
		if pos.prev != nil {
			pos = pos.prev
		} else {
			pos.prev = &n
			n.next = pos
			tl.first = &n
			return true
		}
	}
	pos.next = &n
	n.prev = pos
	if tl.last.next != nil {
		tl.last = &n
	}
	return true
}

func (tl *timedList[T]) DeleteExpired() {
	tl.rwm.Lock()
	defer tl.rwm.Unlock()
	if tl.len.Load() == 0 {
		return
	}
	t := time.Now().Add(-tl.storeDuration)
	pos := tl.first
	for {
		if pos.timestamp.After(t) {
			break
		}
		if pos.next == nil {
			tl.first = nil
			tl.last = nil
			tl.len.Store(0)
			break
		}
		pos.next.prev = nil
		tl.first = pos.next
		pos = pos.next
		tl.len.Add(-1)
	}
}

func (tl *timedList[T]) GetAfter(t time.Time) (data []T, oldest time.Time) {
	tl.rwm.RLock()
	defer tl.rwm.RUnlock()
	if tl.len.Load() == 0 {
		return
	}
	pos := tl.last
	for {
		if pos.timestamp.Before(t) {
			break
		}
		data = append(data, pos.data)
		oldest = pos.timestamp
		if pos.prev != nil {
			pos = pos.prev
		} else {
			break
		}
	}
	return
}

func (tl *timedList[T]) GetLast(t time.Duration) ([]T, time.Time) {
	return tl.GetAfter(time.Now().Add(-t))
}

func (tl *timedList[T]) StoredInterval() time.Duration {
	tl.rwm.RLock()
	defer tl.rwm.RUnlock()
	if tl.len.Load() == 0 {
		return 0
	}
	return tl.last.timestamp.Sub(tl.first.timestamp)
}
