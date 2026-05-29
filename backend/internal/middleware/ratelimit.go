// Package middleware provides rate limiting for HTTP endpoints.
//
// rateLimiter is an in-memory sliding-window rate limiter keyed by IP address.
// It tracks request timestamps per IP and rejects requests that exceed the limit
// within the configured time window. A background goroutine periodically purges
// stale entries to prevent unbounded memory growth.
//
// The current configuration limits auth endpoints (login, register) to 5 requests
// per minute per IP. This is sufficient to prevent brute-force attacks while not
// inconveniencing legitimate users.
//
// This is a single-instance in-process limiter. For horizontally-scaled deployments
// it should be replaced with a Redis-backed implementation (Redis is already
// connected as a dependency but not yet used for this purpose).
package middleware

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/narinsak-u/backend/pkg/utils"
)

// rateLimiter implements a sliding-window rate limit using an in-memory map.
// Each key (IP address) maps to a slice of request timestamps. When a request
// arrives, timestamps outside the current window are pruned, then the remaining
// count is compared against the limit. If the limit is exceeded the request is
// rejected; otherwise the current time is appended and the request proceeds.
//
// Fields:
//   - mu: protects the requests map from concurrent reads/writes
//   - requests: maps client IPs to their recent request timestamps
//   - limit: max requests allowed within the time window
//   - window: duration of the sliding time window
type rateLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

// newRateLimiter creates a rateLimiter and starts a background cleanup goroutine
// that runs every 60 seconds to evict IPs with no activity within the window.
func newRateLimiter(limit int, window time.Duration) *rateLimiter {
	rl := &rateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
	go rl.cleanup()
	return rl
}

// cleanup runs on a 60-second tick and removes IP entries whose timestamps are
// all outside the window. This prevents the map from growing unboundedly with
// IPs that made one request long ago and never returned.
func (rl *rateLimiter) cleanup() {
	for range time.Tick(time.Minute) {
		rl.mu.Lock()
		for ip, times := range rl.requests {
			cutoff := time.Now().Add(-rl.window)
			valid := make([]time.Time, 0)
			for _, t := range times {
				if t.After(cutoff) {
					valid = append(valid, t)
				}
			}
			if len(valid) == 0 {
				delete(rl.requests, ip)
			} else {
				rl.requests[ip] = valid
			}
		}
		rl.mu.Unlock()
	}
}

// allow checks whether a request from the given IP should be allowed.
// It returns true if the IP is under the rate limit, false if it has exceeded it.
// The sliding-window logic:
//  1. Prune timestamps older than the window
//  2. If remaining count >= limit → reject (return false)
//  3. Otherwise → record the current timestamp and allow (return true)
func (rl *rateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	// Filter out expired timestamps
	times := rl.requests[ip]
	valid := make([]time.Time, 0)
	for _, t := range times {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}

	// Compare against the limit
	if len(valid) >= rl.limit {
		rl.requests[ip] = valid
		return false
	}

	// Record this request and allow it
	rl.requests[ip] = append(valid, now)
	return true
}

// authRateLimiter is the singleton instance applied to authentication endpoints.
// It allows 5 requests per minute per IP. This rate is chosen to:
//   - Block automated brute-force password guessing
//   - Allow normal users a few retries on typo'd passwords
//   - Trigger after about 3 minutes of sustained rapid-fire attempts
var authRateLimiter = newRateLimiter(5, time.Minute)

// RateLimitAuth returns a Fiber middleware handler that applies the auth rate
// limiter to incoming requests based on the client's IP address. When the limit
// is exceeded it returns a 429 Too Many Requests response instead of passing
// the request to the next handler.
//
// Usage (in main.go):
//
//	authGroup.Post("/login", middleware.RateLimitAuth(), auth.Login)
func RateLimitAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ip := c.IP()
		if !authRateLimiter.allow(ip) {
			return utils.Error(c, 429, "too many requests. Please try again later")
		}
		return c.Next()
	}
}
