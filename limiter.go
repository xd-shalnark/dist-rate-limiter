package main

import (
	"sync"
	"time"
)

type SlidingWindowLimiter struct {
	mu          sync.Mutex
	prevCount   int
	currCount   int
	limit       int
	windowSize  time.Duration
	windowStart time.Time
}

func (l *SlidingWindowLimiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	elapsed := now.Sub(l.windowStart)
	if elapsed >= l.windowSize {
		l.prevCount = l.currCount
		l.currCount = 0
		l.windowStart = now
	}
	progress := float64(elapsed) / float64(l.windowSize)
	prevWeight := 1.0 - progress
	expectedCount := float64(l.prevCount)*prevWeight + float64(l.currCount)
	if expectedCount >= float64(l.limit) {
		return false
	}
	l.currCount++
	return true
}

func main() {
	limiter := &SlidingWindowLimiter{
		windowSize:  1 * time.Minute,
		limit:       5,
		windowStart: time.Now(),
		prevCount:   0,
		currCount:   0,
	}

	for i := 0; i < 10; i++ {
		if limiter.Allow() {
			println("Request allowed")
		} else {
			println("Request denied")
		}
		time.Sleep(1 * time.Second)
	}
}
