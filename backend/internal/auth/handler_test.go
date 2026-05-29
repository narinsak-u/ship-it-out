package auth

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/narinsak-u/backend/internal/config"
	"github.com/narinsak-u/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func init() {
	config.App = config.Config{
		JWTSecret: "test-secret",
		JWTTTL:    24 * time.Minute,
		Port:      "8080",
	}
}

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *mockRepo) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *mockRepo) FindByID(id uint) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func TestRegister_Success(t *testing.T) {
	repo := new(mockRepo)
	repo.On("Create", mock.MatchedBy(func(u *models.User) bool {
		return u.Name == "John" && u.Email == "john@test.com"
	})).Return(nil)

	app := fiber.New()
	h := NewHandler(repo)
	app.Post("/api/auth/register", h.Register)

	body := `{"name":"John","email":"john@test.com","password":"pass123"}`
	req := httptest.NewRequest("POST", "/api/auth/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, 1000)

	assert.Equal(t, 200, resp.StatusCode)
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.True(t, result["success"].(bool))
	repo.AssertExpectations(t)
}

func TestRegister_ValidationErrors(t *testing.T) {
	repo := new(mockRepo)
	app := fiber.New()
	h := NewHandler(repo)
	app.Post("/api/auth/register", h.Register)

	for _, tt := range []struct{ name, body string }{
		{"missing name", `{"email":"e@t.com","password":"p123"}`},
		{"missing email", `{"name":"John","password":"p123"}`},
		{"missing pass", `{"name":"John","email":"e@t.com"}`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/api/auth/register", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := app.Test(req, 1000)
			assert.Equal(t, 400, resp.StatusCode)
		})
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	repo := new(mockRepo)
	repo.On("Create", mock.Anything).Return(gorm.ErrDuplicatedKey)

	app := fiber.New()
	h := NewHandler(repo)
	app.Post("/api/auth/register", h.Register)

	body := `{"name":"John","email":"dup@test.com","password":"pass123"}`
	req := httptest.NewRequest("POST", "/api/auth/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, 1000)
	assert.Equal(t, 409, resp.StatusCode)
}

func TestLogin_Success(t *testing.T) {
	repo := new(mockRepo)
	hash, _ := bcrypt.GenerateFromPassword([]byte("pass123"), bcrypt.DefaultCost)
	repo.On("FindByEmail", "john@test.com").Return(&models.User{
		ID: 1, Name: "John", Email: "john@test.com",
		Password: string(hash), Role: "customer",
	}, nil)

	app := fiber.New()
	h := NewHandler(repo)
	app.Post("/api/auth/login", h.Login)

	body := `{"email":"john@test.com","password":"pass123"}`
	req := httptest.NewRequest("POST", "/api/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, 1000)

	assert.Equal(t, 200, resp.StatusCode)
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	data := result["data"].(map[string]interface{})
	assert.NotEmpty(t, data["token"])
	repo.AssertExpectations(t)
}

func TestLogin_WrongPassword(t *testing.T) {
	repo := new(mockRepo)
	hash, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)
	repo.On("FindByEmail", "john@test.com").Return(&models.User{
		ID: 1, Email: "john@test.com", Password: string(hash),
	}, nil)

	app := fiber.New()
	h := NewHandler(repo)
	app.Post("/api/auth/login", h.Login)

	body := `{"email":"john@test.com","password":"wrong"}`
	req := httptest.NewRequest("POST", "/api/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, 1000)
	assert.Equal(t, 401, resp.StatusCode)
}

func TestLogin_UserNotFound(t *testing.T) {
	repo := new(mockRepo)
	repo.On("FindByEmail", "unknown@test.com").Return(nil, gorm.ErrRecordNotFound)

	app := fiber.New()
	h := NewHandler(repo)
	app.Post("/api/auth/login", h.Login)

	body := `{"email":"unknown@test.com","password":"pass123"}`
	req := httptest.NewRequest("POST", "/api/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, 1000)
	assert.Equal(t, 401, resp.StatusCode)
}

func TestMe_Success(t *testing.T) {
	repo := new(mockRepo)
	repo.On("FindByID", uint(1)).Return(&models.User{
		ID: 1, Name: "John", Email: "john@test.com", Role: "customer",
	}, nil)

	app := fiber.New()
	h := NewHandler(repo)
	app.Get("/api/auth/me", func(c *fiber.Ctx) error {
		c.Locals("user_id", uint(1))
		return h.Me(c)
	})

	resp, _ := app.Test(httptest.NewRequest("GET", "/api/auth/me", nil), 1000)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestMe_NotAuthenticated(t *testing.T) {
	repo := new(mockRepo)
	app := fiber.New()
	h := NewHandler(repo)
	app.Get("/api/auth/me", h.Me)

	resp, _ := app.Test(httptest.NewRequest("GET", "/api/auth/me", nil), 1000)
	assert.Equal(t, 401, resp.StatusCode)
}

func TestMe_UserNotFound(t *testing.T) {
	repo := new(mockRepo)
	repo.On("FindByID", uint(999)).Return(nil, gorm.ErrRecordNotFound)

	app := fiber.New()
	h := NewHandler(repo)
	app.Get("/api/auth/me", func(c *fiber.Ctx) error {
		c.Locals("user_id", uint(999))
		return h.Me(c)
	})

	resp, _ := app.Test(httptest.NewRequest("GET", "/api/auth/me", nil), 1000)
	assert.Equal(t, 404, resp.StatusCode)
}

func TestLogout(t *testing.T) {
	repo := new(mockRepo)
	app := fiber.New()
	h := NewHandler(repo)
	app.Post("/api/auth/logout", h.Logout)

	resp, _ := app.Test(httptest.NewRequest("POST", "/api/auth/logout", nil), 1000)
	assert.Equal(t, 200, resp.StatusCode)
}
