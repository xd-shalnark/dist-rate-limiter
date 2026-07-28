package registry

import (
	"sync"
	"time"

	"ratelimiter/limiter"
)

// RateLimiter manages a separate sliding-window limiter per key
type RateLimiter struct {
	mu       sync.RWMutex
	limiters map[string]*limiter.SlidingWindowLimiter
	limit    int
	window   time.Duration
}

// Allow reports whether a request for the given key is allowed.
// Creates a new limiter for the key on first use.
func (r *RateLimiter) Allow(key string) bool {
	r.mu.Lock()
	lim, exists := r.limiters[key]
	if !exists {
		lim = limiter.NewSlidingWindowLimiter(r.limit, r.window)
		r.limiters[key] = lim
	}
	r.mu.Unlock() // unlock early, Allow() itself is thread-safe

	return lim.Allow()
}

// NewRateLimiter creates a RateLimiter allowing up
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*limiter.SlidingWindowLimiter),
		limit:    limit,
		window:   window,
	}
}
