package clients

import (
	"sync"
	"time"
)

type RateLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
	limit    int
	interval time.Duration
}

func NewRateLimiter(limit int, interval time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		interval: interval,
	}
}

func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	requests := rl.requests[key]

	// Filter out requests that are outside of the interval
	newRequests := requests[:0] // Use the existing slice to avoid reallocation
	for _, t := range requests {
		if now.Sub(t) < rl.interval {
			newRequests = append(newRequests, t)
		}
	}

	rl.requests[key] = newRequests

	if len(newRequests) >= rl.limit {
		return false
	}

	// Add the current request time to the list
	rl.requests[key] = append(rl.requests[key], now)
	return true
}
