package limiter

import ( // Importing necessary packages
	"sync"
	"time"
)

type SlidingWindowLimiter struct {
	mu          sync.Mutex
	prevCount   int
	currCount   int //Limiter structure
	limit       int
	windowSize  time.Duration
	windowStart time.Time // start time of the current window
}

func (l *SlidingWindowLimiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock() // Locking to ensure thread safety

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
	prevWeight := 1.0 - progress //Mathematical part
	expectedCount := float64(l.prevCount)*prevWeight + float64(l.currCount)

	if expectedCount >= float64(l.limit) {
		return false
	} // allow/deny
	l.currCount++
	return true
}

func NewSlidingWindowLimiter(limit int, windowSize time.Duration) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		limit:       limit,
		windowSize:  windowSize,
		windowStart: time.Now(),
	}
}
