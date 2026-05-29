package middleware

import (
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestRateLimiter_AllowsUpToLimit(t *testing.T) {
	rl := newRateLimiter(3, time.Minute)
	ip := "192.168.1.1"
	assert.True(t, rl.allow(ip))
	assert.True(t, rl.allow(ip))
	assert.True(t, rl.allow(ip))
	assert.False(t, rl.allow(ip))
}

func TestRateLimiter_DifferentIPsAreIndependent(t *testing.T) {
	rl := newRateLimiter(2, time.Minute)
	assert.True(t, rl.allow("10.0.0.1"))
	assert.True(t, rl.allow("10.0.0.1"))
	assert.False(t, rl.allow("10.0.0.1"))
	assert.True(t, rl.allow("10.0.0.2"))
	assert.True(t, rl.allow("10.0.0.2"))
	assert.False(t, rl.allow("10.0.0.2"))
}

func TestRateLimiter_WindowExpiry(t *testing.T) {
	rl := &rateLimiter{
		requests: make(map[string][]time.Time),
		limit:    1,
		window:   30 * time.Millisecond,
	}
	ip := "192.168.1.1"
	assert.True(t, rl.allow(ip))
	assert.False(t, rl.allow(ip))
	time.Sleep(40 * time.Millisecond)
	assert.True(t, rl.allow(ip))
}

func TestRateLimitAuth_Handler(t *testing.T) {
	orig := authRateLimiter
	authRateLimiter = newRateLimiter(2, time.Minute)
	defer func() { authRateLimiter = orig }()

	app := fiber.New()
	app.Post("/test", RateLimitAuth(), func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	resp1, _ := app.Test(httptest.NewRequest("POST", "/test", nil), 1000)
	assert.Equal(t, 200, resp1.StatusCode)

	resp2, _ := app.Test(httptest.NewRequest("POST", "/test", nil), 1000)
	assert.Equal(t, 200, resp2.StatusCode)

	resp3, _ := app.Test(httptest.NewRequest("POST", "/test", nil), 1000)
	assert.Equal(t, 429, resp3.StatusCode)
}

func TestRateLimiter_ConcurrentAccess(t *testing.T) {
	rl := newRateLimiter(100, time.Minute)
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			rl.allow("192.168.1.1")
		}()
	}
	wg.Wait()
}
