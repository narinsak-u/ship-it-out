package presentation

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/narinsak-u/backend/internal/auth/domain"
	"github.com/narinsak-u/backend/internal/platform/config"
)

type currentUserStub struct{ user *domain.User }

func (r *currentUserStub) FindByID(uint) (*domain.User, error) { return r.user, nil }

func TestStaffOnlyChecksPersistedRole(t *testing.T) {
	oldConfig := config.App
	config.App.JWTSecret = "test-secret"
	t.Cleanup(func() { config.App = oldConfig })

	app := fiber.New()
	user := &domain.User{ID: 7, Role: domain.RoleStaff}
	app.Post("/write", AuthRequired(), StaffOnly(&currentUserStub{user: user}), func(c *fiber.Ctx) error {
		return c.SendStatus(204)
	})
	makeToken := func(role string) string {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": 7,
			"role":    role,
			"exp":     time.Now().Add(time.Hour).Unix(),
		})
		value, err := token.SignedString([]byte(config.App.JWTSecret))
		if err != nil {
			t.Fatal(err)
		}
		return value
	}
	request := func(token string) int {
		req := httptest.NewRequest("POST", "/write", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		response, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		response.Body.Close()
		return response.StatusCode
	}

	if got := request(makeToken(domain.RoleCustomer)); got != 204 {
		t.Fatalf("staff database role with customer claim status = %d, want 204", got)
	}
	user.Role = domain.RoleCustomer
	if got := request(makeToken(domain.RoleStaff)); got != 403 {
		t.Fatalf("downgraded staff token status = %d, want 403", got)
	}
}
