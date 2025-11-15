package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"
)

// rateLimiter keeps track of requests per IP in a sliding window.
type rateLimiter struct {
	mu     sync.Mutex             // A lock so two requests don’t modify data at the same time. Important for concurrency.
	limit  int                    //Maximum allowed requests
	window time.Duration          //The time window (example: 10 minute).
	visits map[string][]time.Time // Map of: IP to list of timestamps for each request.
}

//visits["127.0.0.1"] = [2025-11-15 18:00:01, 18:00:05, 18:00:10]

// Constructor creates a new rate limiter instance.
func NewRateLimiter(limit int, window time.Duration) *rateLimiter {
	return &rateLimiter{
		limit:  limit,
		window: window,
		visits: make(map[string][]time.Time),
	}
}

// allow registers a request from an IP and returns true if allowed.
func (rl *rateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	// Filter timestamps inside the window
	times := rl.visits[ip]
	var filtered []time.Time
	for _, t := range times {
		if t.After(cutoff) {
			filtered = append(filtered, t)
		}
	}

	// Check limit
	if len(filtered) >= rl.limit {
		rl.visits[ip] = filtered
		return false // block request
	}

	// Add this visit
	filtered = append(filtered, now)
	rl.visits[ip] = filtered
	return true
}

// Middleware wraps an existing handler with rate limiting.
func (rl *rateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := clientIP(r)
		if !rl.allow(ip) {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("rate limit exceeded, try again later"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Extract client IP from request.
func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
