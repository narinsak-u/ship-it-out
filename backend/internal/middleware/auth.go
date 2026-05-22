package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/narinsak-u/backend/internal/config"
	"github.com/narinsak-u/backend/pkg/utils"
)

func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		tokenStr := ""

		if auth != "" && strings.HasPrefix(auth, "Bearer ") {
			tokenStr = strings.TrimPrefix(auth, "Bearer ")
		} else {
			// Fallback to jwt cookie
			tokenStr = c.Cookies("jwt")
		}

		if tokenStr == "" {
			return utils.Error(c, 401, "missing or invalid token")
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(config.App.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			return utils.Error(c, 401, "invalid or expired token")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return utils.Error(c, 401, "invalid token claims")
		}

		uid, _ := claims["user_id"].(float64)
		c.Locals("user_id", uint(uid))
		c.Locals("role", claims["role"])
		return c.Next()
	}
}
