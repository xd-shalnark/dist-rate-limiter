package registry

import (
	"sync"
	"time"

	"ratelimiter/limiter"
)

type RateLimiter struct {
	mu       sync.RWMutex
	limiters map[string]*limiter.SlidingWindowLimiter
	limit    int
	window   time.Duration
}

func (r *RateLimiter) Allow(key string) bool {
	r.mu.Lock()
	lim, exists := r.limiters[key]
	if !exists {
		lim = limiter.NewSlidingWindowLimiter(r.limit, r.window)
		r.limiters[key] = lim
	}
	r.mu.Unlock()

	return lim.Allow()
}
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*limiter.SlidingWindowLimiter),
		limit:    limit,
		window:   window,
	}
}