package registry

import (
	"fmt"
	"sync"
	"time"

	"ratelimiter/limiter"
)

type entry struct {
	lim      *limiter.SlidingWindowLimiter
	lastSeen time.Time
}

type RateLimiter struct {
	mu       sync.RWMutex
	limiters map[string]*entry
	limit    int
	window   time.Duration
	ttl      time.Duration
}

func (r *RateLimiter) Allow(key string) bool {
	r.mu.Lock()
	e, exists := r.limiters[key]
	if !exists {

		lim, err := limiter.NewSlidingWindowLimiter(r.limit, r.window)
		if err != nil {
			panic(fmt.Sprintf("registry: invariant violated, limit/window became invalid after construction: %v", err))
		}
		e = &entry{lim: lim}
		r.limiters[key] = e
	}
	e.lastSeen = time.Now()
	r.mu.Unlock()

	return e.lim.Allow()
}

func (r *RateLimiter) cleanupLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		r.mu.Lock()
		for k, e := range r.limiters {
			if time.Since(e.lastSeen) > r.ttl {
				delete(r.limiters, k)
			}
		}
		r.mu.Unlock()
	}
}

// NewRateLimiter creates a RateLimiter allowing up to `limit` requests
// per `window`, evicting idle keys after `ttl` of inactivity.
// Returns an error if limit/window are invalid — caught once at startup
// instead of panicking later on the first request.
func NewRateLimiter(limit int, window, ttl time.Duration) (*RateLimiter, error) {

	if _, err := limiter.NewSlidingWindowLimiter(limit, window); err != nil {
		return nil, fmt.Errorf("registry: invalid RateLimiter config: %w", err)
	}
	if ttl <= 0 {
		return nil, fmt.Errorf("registry: ttl must be positive, got %v", ttl)
	}

	r := &RateLimiter{
		limiters: make(map[string]*entry),
		limit:    limit,
		window:   window,
		ttl:      ttl,
	}
	go r.cleanupLoop(time.Minute)
	return r, nil
}
