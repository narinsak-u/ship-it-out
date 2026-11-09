package auth

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/narinsak-u/backend/internal/platform/config"
	"github.com/narinsak-u/backend/internal/models"
	"github.com/narinsak-u/backend/pkg/utils"
)

// Handler parses auth requests and formats HTTP responses around Service workflows.
type Handler struct {
	service *Service
}

// NewHandler creates an HTTP handler backed by the auth service.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// The name of the cookie where we store the JWT token for the browser.
const cookieName = "jwt"

// RegisterRequest is the JSON body the client sends when creating a new account.
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// LoginRequest is the JSON body the client sends when signing in.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// safeUser contains the user fields exposed by auth endpoints.
type safeUser struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// userResponse prevents credentials from escaping through user responses.
func userResponse(user *models.User) safeUser {
	return safeUser{ID: user.ID, Name: user.Name, Email: user.Email, Role: user.Role, CreatedAt: user.CreatedAt}
}

// setAuthCookie stores the signed JWT in the browser cookie.
func setAuthCookie(c *fiber.Ctx, token string) {
	c.Cookie(&fiber.Cookie{
		Name:     cookieName,
		Value:    token,
		Path:     "/",
		HTTPOnly: true,
		SameSite: "None",
		Secure:   true,
		MaxAge:   int(config.App.JWTTTL.Seconds()),
	})
}

// Register creates a read-only customer account and starts a session.
func (h *Handler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, 400, "invalid request body")
	}
	if req.Name == "" || req.Email == "" || req.Password == "" {
		return utils.Error(c, 400, "name, email, and password are required")
	}
	user, token, err := h.service.Register(req.Name, req.Email, req.Password)
	if errors.Is(err, ErrEmailAlreadyRegistered) {
		return utils.Error(c, 409, "email already registered")
	}
	if errors.Is(err, ErrPasswordHashing) {
		return utils.Error(c, 500, "failed to hash password")
	}
	if errors.Is(err, ErrTokenGeneration) {
		return utils.Error(c, 500, "failed to generate token")
	}
	if err != nil {
		return utils.Error(c, 500, "failed to register user")
	}
	setAuthCookie(c, token.String())
	return utils.Success(c, fiber.Map{"user": userResponse(user)})
}

// Login authenticates an existing user and returns its token and safe profile.
func (h *Handler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, 400, "invalid request body")
	}
	user, token, err := h.service.Login(req.Email, req.Password)
	if errors.Is(err, ErrInvalidCredentials) {
		return utils.Error(c, 401, "invalid email or password")
	}
	if errors.Is(err, ErrTokenGeneration) {
		return utils.Error(c, 500, "failed to generate token")
	}
	if err != nil {
		return utils.Error(c, 500, "failed to authenticate")
	}
	tokenString := token.String()
	setAuthCookie(c, tokenString)
	return utils.Success(c, fiber.Map{"token": tokenString, "user": userResponse(user)})
}

// Me returns the currently authenticated user's profile.
func (h *Handler) Me(c *fiber.Ctx) error {
	// AuthRequired middleware sets user_id after verifying the JWT.
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.Error(c, 401, "not authenticated")
	}
	user, err := h.service.Profile(userID)
	if errors.Is(err, ErrUserNotFound) {
		return utils.Error(c, 404, "user not found")
	}
	if err != nil {
		return utils.Error(c, 500, "failed to load user profile")
	}
	return utils.Success(c, userResponse(user))
}

// Logout clears the auth cookie so the user is no longer authenticated.
func (h *Handler) Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		HTTPOnly: true,
		SameSite: "None",
		Secure:   true,
		MaxAge:   0,
	})
	return utils.Success(c, fiber.Map{"message": "logged out"})
}
