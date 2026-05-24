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

// The name of the cookie where we store the JWT token for the browser
const cookieName = "jwt"

// RegisterRequest is the JSON body the client sends when creating a new account
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// LoginRequest is the JSON body the client sends when signing in
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// setAuthCookie writes the JWT token into an HTTP-only cookie on the response.
// The cookie is:
//   - HTTP-only (JavaScript can't read it — prevents XSS attacks)
//   - SameSite=Lax (sent on same-site navigation but not cross-site)
//   - MaxAge=86400 (24 hours — after that the browser deletes it)
//   - Path="/" (sent to every page on this domain)
func setAuthCookie(c *fiber.Ctx, token string) {
	c.Cookie(&fiber.Cookie{
	role := "customer"
		Value:    token,
		Path:     "/",
		HTTPOnly: true,
		SameSite: "Lax",
		MaxAge:   86400,
	})
}

// Register creates a new user account. Steps:
//  1. Read the JSON body (name, email, password, optional role) from the request
//  2. Validate that required fields are present
//  3. Hash the password so we never store it in plain text
//  4. If no role was provided, default to "customer"
//  5. Save the new user to the database (fails with 409 if email already exists)
//  6. Create a signed JWT token containing user_id, role, and expiry time
//  7. Set the JWT as an HTTP-only cookie so the browser remembers the session
//  8. Return the new user object as JSON
func Register(c *fiber.Ctx) error {
	// --- Parse + validate the incoming JSON ---
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, 400, "invalid request body")
	}
	if req.Name == "" || req.Email == "" || req.Password == "" {
		return utils.Error(c, 400, "name, email, and password are required")
	}

	// --- Hash the password before storing it ---
	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return utils.Error(c, 500, "failed to hash password")
	}

	// --- Default to "customer" role if not specified ---
	role := req.Role
	if role == "" {
		role = "customer"
	}

	// --- Build the user model and insert into Postgres ---
	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hash,
		Role:     role,
	}

	// If the email is already taken, the unique constraint fails → 409 Conflict
	if result := database.DB.Create(&user); result.Error != nil {
		return utils.Error(c, 409, "email already registered")
	}

	// --- Issue a JWT so the user is logged in right after registering ---
	// Claims = data embedded inside the token (signed, not encrypted)
	claims := jwt.MapClaims{
		"user_id": user.ID,                                  // who the user is
		"role":    user.Role,                                // what they can do
		"exp":     time.Now().Add(config.App.JWTTTL).Unix(), // when it expires (Unix timestamp)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(config.App.JWTSecret))
	if err != nil {
		return utils.Error(c, 500, "failed to generate token")
	}

	// --- Store the token in a cookie the browser will send on subsequent requests ---
	setAuthCookie(c, tokenStr)

	return utils.Success(c, fiber.Map{"user": user})
}

// Login authenticates an existing user. Steps:
//  1. Read the JSON body (email, password) from the request
//  2. Look up the user by email in the database (401 if not found)
//  3. Compare the provided password against the stored hash (401 if wrong)
//  4. Create a signed JWT token containing user_id, role, and expiry time
//  5. Set the JWT as an HTTP-only cookie
//  6. Return the token + user object as JSON (so mobile/CLI clients can use the token directly)
func Login(c *fiber.Ctx) error {
	// --- Parse the incoming JSON ---
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, 400, "invalid request body")
	}

	// --- Find the user by email ---
	var user models.User
	if result := database.DB.Where("email = ?", req.Email).First(&user); result.Error != nil {
		// We say "invalid email or password" rather than "user not found"
		// to avoid leaking which emails are registered
		return utils.Error(c, 401, "invalid email or password")
	}

	// --- Verify the password against the stored hash ---
	if !utils.CheckPassword(req.Password, user.Password) {
		return utils.Error(c, 401, "invalid email or password")
	}

	// --- Issue a JWT for this session ---
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

	// --- Set the cookie AND return the token in the JSON body ---
	setAuthCookie(c, tokenStr)

	return utils.Success(c, fiber.Map{
		"token": tokenStr,
		"user":  user,
	})
}

// Me returns the currently authenticated user's profile.
// It relies on the AuthRequired middleware having already verified the JWT
// and placed the user_id into c.Locals("user_id").
// Steps:
//  1. Read user_id from the request context (set by middleware)
//  2. Look up the user in the database by that ID
//  3. Return a safe subset of user fields (no password hash!)
func Me(c *fiber.Ctx) error {
	// AuthRequired middleware ran before us and stored user_id in locals
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.Error(c, 401, "not authenticated")
	}

	// Fetch fresh user data from the database
	var user models.User
	if result := database.DB.First(&user, userID); result.Error != nil {
		return utils.Error(c, 404, "user not found")
	}

	// Return only safe fields — never expose the password hash
	return utils.Success(c, fiber.Map{
		"id":         user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"role":       user.Role,
		"created_at": user.CreatedAt,
	})
}

// Logout clears the auth cookie so the user is no longer authenticated.
// Steps:
//  1. Overwrite the "jwt" cookie with an empty value and MaxAge=0
//     (MaxAge=0 tells the browser to delete the cookie immediately)
//  2. Return a success message
func Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		HTTPOnly: true,
		SameSite: "Lax",
		MaxAge:   0, // 0 = delete the cookie right now
	})
	return utils.Success(c, fiber.Map{"message": "logged out"})
}
