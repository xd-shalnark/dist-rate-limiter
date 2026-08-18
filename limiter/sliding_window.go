package limiter

import (
	"fmt"
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
		windowsPassed := int(elapsed / l.windowSize)
		if windowsPassed == 1 {
			l.prevCount = l.currCount
		} else {
			l.prevCount = 0
		}
		l.currCount = 0
		l.windowStart = l.windowStart.Add(time.Duration(windowsPassed) * l.windowSize)
		elapsed = now.Sub(l.windowStart)
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

// "integer divide by zero" в Allow(), если windowSize == 0.
func NewSlidingWindowLimiter(limit int, windowSize time.Duration) (*SlidingWindowLimiter, error) {
	if limit <= 0 {
		return nil, fmt.Errorf("limiter: limit must be positive, got %d", limit)
	}
	if windowSize <= 0 {
		return nil, fmt.Errorf("limiter: windowSize must be positive, got %v", windowSize)
	}
	return &SlidingWindowLimiter{
		limit:       limit,
		windowSize:  windowSize,
		windowStart: time.Now(),
	}, nil
}
