package presentation

import (
	"errors"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gofiber/fiber/v2"
	"github.com/narinsak-u/backend/internal/auth/application"
	"github.com/narinsak-u/backend/internal/auth/domain"
	"github.com/narinsak-u/backend/internal/platform/config"
	utils "github.com/narinsak-u/backend/internal/platform/http"
)

type TokenClaims struct {
	UserID uint
	Role   string
}

var ErrInvalidTokenClaims = errors.New("invalid token claims")

func ParseToken(value string, secret []byte) (TokenClaims, error) {
	token, err := jwt.Parse(value, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})
	if err != nil || !token.Valid {
		return TokenClaims{}, errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return TokenClaims{}, ErrInvalidTokenClaims
	}
	userID, ok := claims["user_id"].(float64)
	if !ok || userID <= 0 {
		return TokenClaims{}, ErrInvalidTokenClaims
	}
	role, ok := claims["role"].(string)
	if !ok {
		return TokenClaims{}, ErrInvalidTokenClaims
	}
	return TokenClaims{UserID: uint(userID), Role: role}, nil
}

type CurrentUserRepository = application.CurrentUserRepository

func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authorization := c.Get("Authorization")
		token := ""
		if authorization != "" && strings.HasPrefix(authorization, "Bearer ") {
			token = strings.TrimPrefix(authorization, "Bearer ")
		} else {
			token = c.Cookies("jwt")
		}
		if token == "" {
			return utils.Error(c, 401, "missing or invalid token")
		}
		claims, err := ParseToken(token, []byte(config.App.JWTSecret))
		if errors.Is(err, ErrInvalidTokenClaims) {
			return utils.Error(c, 401, "invalid token claims")
		}
		if err != nil {
			return utils.Error(c, 401, "invalid or expired token")
		}
		c.Locals("user_id", claims.UserID)
		c.Locals("role", claims.Role)
		return c.Next()
	}
}

func StaffOnly(repo CurrentUserRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := c.Locals("user_id").(uint)
		if !ok {
			return utils.Error(c, 401, "missing or invalid token")
		}
		user, err := repo.FindByID(userID)
		if err != nil {
			return utils.Error(c, 401, "user not found")
		}
		if user.Role != domain.RoleAdmin && user.Role != domain.RoleStaff {
			return utils.Error(c, 403, "staff access required")
		}
		return c.Next()
	}
}
