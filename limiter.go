package main

import (   // Importing necessary packages
	"sync"
	"time"
)

type RateLimiter struct {
    mu       sync.RWMutex
    limiters map[string]*SlidingWindowLimiter 
    limit    int
    window   time.Duration
}

type SlidingWindowLimiter struct {
	mu          sync.Mutex
	prevCount   int
	currCount   int 		//Limiter structure
	limit       int
	windowSize  time.Duration
	windowStart time.Time
}

func (r *RateLimiter) Allow(key string) bool {
    r.mu.Lock()
    limiter, exists := r.limiters[key]
    if !exists {
        limiter = &SlidingWindowLimiter{
            limit:       r.limit,
            windowSize:  r.window,
            windowStart: time.Now(),
        }
        r.limiters[key] = limiter
    }
    r.mu.Unlock()

    return limiter.Allow()
}

func (l *SlidingWindowLimiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock() 		// Locking to ensure thread safety

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
	prevWeight := 1.0 - progress   //Mathematical part
	expectedCount := float64(l.prevCount)*prevWeight + float64(l.currCount)

	if expectedCount >= float64(l.limit) {
		return false
	}                 					// allow/deny
	l.currCount++
	return true
}

func main() {
	rl := &RateLimiter{
		limiters: make(map[string]*SlidingWindowLimiter),
		limit:    3,
		window:   1 * time.Minute,
	}

	for i := 0; i < 5; i++ {
		if rl.Allow("user-A") {
			println("user-A: allowed")
		} else {
			println("user-A: denied")
		}
	}

	if rl.Allow("user-B") {
		println("user-B: allowed")
	} else {
		println("user-B: denied")
	}
}
