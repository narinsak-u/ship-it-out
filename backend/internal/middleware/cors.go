package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/narinsak-u/backend/internal/config"
)

// CORS returns middleware that sets Cross-Origin Resource Sharing (CORS)
// headers on every response. CORS is a browser security mechanism that
// controls which websites are allowed to make requests to this API.
//
// Why this is needed:
//   - The frontend runs on http://localhost:5173 (Vite dev server)
//   - The backend runs on http://localhost:8080 (or whatever port is configured)
//   - Browsers block requests from one origin to another unless CORS headers
//     explicitly allow it
//
// What it does:
//  1. Allows requests from the frontend's origin (localhost:5173)
//  2. Allows credentials (cookies) to be sent cross-origin — needed so the
//     browser sends the "jwt" cookie to the API
//  3. Lists permitted HTTP methods and headers
//  4. Preflight requests (OPTIONS) are answered immediately with 204 No Content
//     so the browser knows the actual request is safe to send
func CORS() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", config.App.CORSOrigin)
		c.Set("Access-Control-Allow-Credentials", "true")
		c.Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		if c.Method() == "OPTIONS" {
			return c.SendStatus(204)
		}
		return c.Next()
	}
}
