// Package middleware provides reusable Fiber middleware that runs before
// route handlers. Each function returns a fiber.Handler that can be plugged
// into any route or route group.
package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/narinsak-u/backend/internal/config"
	"github.com/narinsak-u/backend/pkg/utils"
)

// AuthRequired returns middleware that guards routes behind authentication.
// It checks for a valid JWT in two places (in order):
//  1. The "Authorization: Bearer <token>" header (used by mobile/CLI clients)
//  2. The "jwt" cookie (used by browser-based clients)
//
// If a valid token is found, it extracts the user_id and role from the JWT
// claims and stores them in the request context (c.Locals) so downstream
// handlers can access them (e.g. auth.Me reads c.Locals("user_id")).
//
// If no token is found or the token is invalid/expired, it returns 401 and
// blocks the request from reaching the route handler.
func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// --- Try to get the token from the Authorization header first ---
		auth := c.Get("Authorization")
		tokenStr := ""

		if auth != "" && strings.HasPrefix(auth, "Bearer ") {
			// Extract just the token part after "Bearer "
			tokenStr = strings.TrimPrefix(auth, "Bearer ")
		} else {
			// Fallback: check for the "jwt" cookie (set by auth handlers on login/register)
			tokenStr = c.Cookies("jwt")
		}

		// --- Reject if no token was found in either location ---
		if tokenStr == "" {
			return utils.Error(c, 401, "missing or invalid token")
		}

		// --- Parse and verify the JWT signature ---
		// jwt.Parse() automatically validates the HMAC signature using our secret
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(config.App.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			return utils.Error(c, 401, "invalid or expired token")
		}

		// --- Extract the claims (the data we embedded when creating the token) ---
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return utils.Error(c, 401, "invalid token claims")
		}

		// --- Store user info in the request context for downstream handlers ---
		// JSON numbers decode as float64, so we convert to uint
		uid, _ := claims["user_id"].(float64)
		c.Locals("user_id", uint(uid))
		c.Locals("role", claims["role"])

		return c.Next() // Pass control to the next handler
	}
}
