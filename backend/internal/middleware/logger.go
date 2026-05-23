package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// Logger returns middleware that logs every incoming HTTP request.
//
// How it works:
//  1. Records the current time before the request is processed
//  2. Calls c.Next() to pass control to the route handler (and any subsequent
//     middleware). This blocks until the handler sends a response
//  3. After the response is sent, calculates how long the request took
//  4. Logs the HTTP method, path, response status code, and duration
//
// This produces output like:
//
//	{"method":"GET","path":"/api/shipments","status":200,"duration":"12ms"}
//
// This is useful for debugging, monitoring, and spotting slow endpoints.
func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)
		log.Info().
			Str("method", c.Method()).
			Str("path", c.Path()).
			Int("status", c.Response().StatusCode()).
			Dur("duration", duration).
			Msg("request")
		return err
	}
}
