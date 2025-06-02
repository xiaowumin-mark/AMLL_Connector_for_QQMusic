package main

import (
	"sync"
	"time"
)

type TimerID int

var (
	timerIDCounter TimerID = 0
	timeouts               = make(map[TimerID]*time.Timer)
	intervals              = make(map[TimerID]*time.Ticker)
	muu            sync.Mutex
)

func GsetTimeout(callback func(), delay time.Duration) TimerID {
	muu.Lock()
	defer muu.Unlock()

	timerIDCounter++
	id := timerIDCounter

	timeouts[id] = time.AfterFunc(delay, func() {
		callback()
		muu.Lock()
		delete(timeouts, id)
		muu.Unlock()
	})

	return id
}

func GclearTimeout(id TimerID) {
	muu.Lock()
	defer muu.Unlock()

	if timer, exists := timeouts[id]; exists {
		timer.Stop()
		delete(timeouts, id)
	}
}

func GsetInterval(callback func(), interval time.Duration) TimerID {
	muu.Lock()
	defer muu.Unlock()

	timerIDCounter++
	id := timerIDCounter

	intervals[id] = time.NewTicker(interval)

	go func() {
		for range intervals[id].C {
			callback()
		}
	}()

	return id
}

func GclearInterval(id TimerID) {
	muu.Lock()
	defer muu.Unlock()

	if ticker, exists := intervals[id]; exists {
		ticker.Stop()
		delete(intervals, id)
	}
}

type Throttler struct {
	duration time.Duration
	last     time.Time
	mu       sync.Mutex
}

func NewThrottler(duration time.Duration) *Throttler {
	return &Throttler{duration: duration}
}

// Do 执行节流函数：立即执行首次调用，后续调用在间隔内被忽略
func (t *Throttler) Do(f func()) {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	if t.last.IsZero() || now.Sub(t.last) >= t.duration {
		t.last = now
		f()
	}
}
