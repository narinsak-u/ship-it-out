package auth

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/narinsak-u/backend/internal/config"
	"github.com/narinsak-u/backend/internal/database"
	"github.com/narinsak-u/backend/internal/models"
	"github.com/narinsak-u/backend/pkg/utils"
)

const cookieName = "jwt"

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func setAuthCookie(c *fiber.Ctx, token string) {
	c.Cookie(&fiber.Cookie{
		Name:     cookieName,
		Value:    token,
		Path:     "/",
		HTTPOnly: true,
		SameSite: "Lax",
		MaxAge:   86400,
	})
}

func Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, 400, "invalid request body")
	}
	if req.Name == "" || req.Email == "" || req.Password == "" {
		return utils.Error(c, 400, "name, email, and password are required")
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return utils.Error(c, 500, "failed to hash password")
	}

	role := req.Role
	if role == "" {
		role = "customer"
	}

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hash,
		Role:     role,
	}

	if result := database.DB.Create(&user); result.Error != nil {
		return utils.Error(c, 409, "email already registered")
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(config.App.JWTTTL).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(config.App.JWTSecret))
	if err != nil {
		return utils.Error(c, 500, "failed to generate token")
	}

	setAuthCookie(c, tokenStr)

	return utils.Success(c, fiber.Map{"user": user})
}

func Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, 400, "invalid request body")
	}

	var user models.User
	if result := database.DB.Where("email = ?", req.Email).First(&user); result.Error != nil {
		return utils.Error(c, 401, "invalid email or password")
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return utils.Error(c, 401, "invalid email or password")
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(config.App.JWTTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(config.App.JWTSecret))
	if err != nil {
		return utils.Error(c, 500, "failed to generate token")
	}

	setAuthCookie(c, tokenStr)

	return utils.Success(c, fiber.Map{
		"token": tokenStr,
		"user":  user,
	})
}

func Me(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.Error(c, 401, "not authenticated")
	}

	var user models.User
	if result := database.DB.First(&user, userID); result.Error != nil {
		return utils.Error(c, 404, "user not found")
	}

	return utils.Success(c, fiber.Map{
		"id":         user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"role":       user.Role,
		"created_at": user.CreatedAt,
	})
}

func Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		HTTPOnly: true,
		SameSite: "Lax",
		MaxAge:   0,
	})
	return utils.Success(c, fiber.Map{"message": "logged out"})
}
