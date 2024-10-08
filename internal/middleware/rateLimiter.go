package middleware

import (
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	rateLimit  int
	timeWindow time.Duration
	users      map[string][]time.Time
	mu         sync.Mutex
}

func NewRateLimiter(rateLimit int, timeWindow time.Duration) *RateLimiter {
	return &RateLimiter{
		rateLimit:  rateLimit,
		timeWindow: timeWindow,
		users:      make(map[string][]time.Time),
	}
}

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	requestTimes := rl.users[ip]
	validRequests := make([]time.Time, 0, len(requestTimes))
	for _, t := range requestTimes {
		if now.Sub(t) < rl.timeWindow {
			validRequests = append(validRequests, t)
		}
	}

	if len(validRequests) < rl.rateLimit {
		validRequests = append(validRequests, now)
		rl.users[ip] = validRequests
		return true
	}

	rl.users[ip] = validRequests
	return false
}

func (rl *RateLimiter) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		if !rl.Allow(ip) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
