package handlers

import (
	"net/http"
	"sync"
	"time"
)

type rateLimiter struct {
	maxRequests int
	interval    time.Duration
	mu          sync.Mutex
	counters    map[string][]time.Time
}

func NewRateLimiter(maxRequests int, interval time.Duration) *rateLimiter {
	return &rateLimiter{
		maxRequests: maxRequests,
		interval:    interval,
		counters:    make(map[string][]time.Time),
	}
}

func (rl *rateLimiter) LimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		rl.mu.Lock()
		defer rl.mu.Unlock()
		now := time.Now()
		counter := rl.counters[ip]
		// Remove expired requests from the counter
		for i := len(counter) - 1; i >= 0; i-- {
			if now.Sub(counter[i]) > rl.interval {
				counter = counter[i+1:]
				break
			}
		}
		// Check if the request is allowed
		if len(counter) >= rl.maxRequests {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		// Add the request to the counter
		rl.counters[ip] = append(counter, now)
		next.ServeHTTP(w, r)
	}
}
