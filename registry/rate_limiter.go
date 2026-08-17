package registry

import (
	"sync"
	"time"

	"ratelimiter/limiter"
)

// entry хранит лимитер конкретного ключа + время последнего обращения,
// нужно для TTL-эвикшна в cleanupLoop.
type entry struct {
	lim      *limiter.SlidingWindowLimiter
	lastSeen time.Time
}

// RateLimiter manages a separate sliding-window limiter per key
type RateLimiter struct {
	mu       sync.RWMutex
	limiters map[string]*entry // ← тип поля теперь map[string]*entry, а не *SlidingWindowLimiter
	limit    int
	window   time.Duration
	ttl      time.Duration // ← новое поле, которого не хватало
}

// Allow reports whether a request for the given key is allowed.
// Creates a new limiter for the key on first use.
func (r *RateLimiter) Allow(key string) bool {
	r.mu.Lock()
	e, exists := r.limiters[key]
	if !exists {
		e = &entry{lim: limiter.NewSlidingWindowLimiter(r.limit, r.window)}
		r.limiters[key] = e
	}
	e.lastSeen = time.Now()
	r.mu.Unlock()

	return e.lim.Allow()
}

// cleanupLoop периодически удаляет записи, к которым давно не обращались,
// чтобы карта limiters не росла бесконечно (защита от утечки памяти).
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
// per `window`, and evicts idle keys after `ttl` of inactivity.
func NewRateLimiter(limit int, window, ttl time.Duration) *RateLimiter {
	r := &RateLimiter{
		limiters: make(map[string]*entry),
		limit:    limit,
		window:   window,
		ttl:      ttl,
	}
	go r.cleanupLoop(time.Minute) // фоновая горутина очистки
	return r
}
