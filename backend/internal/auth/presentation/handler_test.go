package presentation

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/narinsak-u/backend/internal/auth/application"
	"github.com/narinsak-u/backend/internal/auth/domain"
)

type workflowsStub struct {
	user *domain.User
	registerErr error
	loginErr error
	profileErr error
}

func (s workflowsStub) Register(string, string, string) (*domain.User, string, error) {
	return s.user, "signed-token", s.registerErr
}
func (s workflowsStub) Login(string, string) (*domain.User, string, error) {
	return s.user, "signed-token", s.loginErr
}
func (s workflowsStub) Profile(uint) (*domain.User, error) { return s.user, s.profileErr }

func authApp(method, path string, handler fiber.Handler) *fiber.App {
	app := fiber.New()
	if method == "GET" {
		app.Get(path, func(c *fiber.Ctx) error { c.Locals("user_id", uint(3)); return handler(c) })
	} else {
		app.Post(path, handler)
	}
	return app
}

func authRequest(t *testing.T, app *fiber.App, method, path, body string) *http.Response {
	t.Helper()
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	response, err := app.Test(request)
	if err != nil { t.Fatal(err) }
	return response
}

func TestRegisterIgnoresRequestedRoleAndHidesPassword(t *testing.T) {
	handler := NewHandler(workflowsStub{user: &domain.User{ID: 3, Name: "N", Email: "n@example.com", Password: "secret-hash", Role: domain.RoleCustomer}})
	response := authRequest(t, authApp("POST", "/api/auth/register", handler.Register), "POST", "/api/auth/register", `{"name":"N","email":"n@example.com","password":"pw","role":"admin"}`)
	defer response.Body.Close()
	var body map[string]any
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil { t.Fatal(err) }
	data := body["data"].(map[string]any)
	user := data["user"].(map[string]any)
	if user["role"] != "customer" { t.Fatalf("role = %v", user["role"]) }
	if _, ok := user["password"]; ok { t.Fatal("password exposed") }
	if response.Header.Get("Set-Cookie") == "" { t.Fatal("auth cookie missing") }
}

func TestLoginReturnsUserAndCookie(t *testing.T) {
	handler := NewHandler(workflowsStub{user: &domain.User{ID: 3, Name: "N", Email: "n@example.com", Password: "secret-hash", Role: domain.RoleCustomer}})
	response := authRequest(t, authApp("POST", "/api/auth/login", handler.Login), "POST", "/api/auth/login", `{"email":"n@example.com","password":"pw"}`)
	defer response.Body.Close()
	var body struct { Success bool `json:"success"`; Data struct { Token string `json:"token"`; User map[string]any `json:"user"` } `json:"data"` }
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil { t.Fatal(err) }
	cookie := response.Cookies()
	if response.StatusCode != 200 || !body.Success || body.Data.Token != "signed-token" || body.Data.User["email"] != "n@example.com" || body.Data.User["password"] != nil || len(cookie) != 1 || cookie[0].Name != "jwt" || !cookie[0].HttpOnly {
		t.Fatalf("login response = %d %+v cookie=%+v", response.StatusCode, body, cookie)
	}
}

func TestLoginRejectsCredentialsAndMalformedRequests(t *testing.T) {
	handler := NewHandler(workflowsStub{loginErr: application.ErrInvalidCredentials})
	for _, tc := range []struct { body, want string; status int }{
		{`{"email":"n@example.com","password":"wrong"}`, "invalid email or password", 401},
		{`{`, "invalid request body", 400},
	} {
		response := authRequest(t, authApp("POST", "/api/auth/login", handler.Login), "POST", "/api/auth/login", tc.body)
		var body map[string]any
		_ = json.NewDecoder(response.Body).Decode(&body)
		response.Body.Close()
		if response.StatusCode != tc.status || body["error"] != tc.want { t.Errorf("login response = %d %v, want %d %q", response.StatusCode, body, tc.status, tc.want) }
	}
}

func TestProfileReturnsSafeUserAndMapsFailures(t *testing.T) {
	for _, tc := range []struct { name string; stub workflowsStub; wantStatus int; wantError string }{
		{"success", workflowsStub{user: &domain.User{ID: 3, Name: "N", Email: "n@example.com", Password: "secret-hash", Role: "staff"}}, 200, ""},
		{"not found", workflowsStub{profileErr: application.ErrUserNotFound}, 404, "user not found"},
		{"failure", workflowsStub{profileErr: errors.New("database unavailable")}, 500, "failed to load user profile"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			response := authRequest(t, authApp("GET", "/api/auth/me", NewHandler(tc.stub).Me), "GET", "/api/auth/me", "")
			defer response.Body.Close()
			var body struct { Error string `json:"error"`; Data map[string]any `json:"data"` }
			if err := json.NewDecoder(response.Body).Decode(&body); err != nil { t.Fatal(err) }
			if response.StatusCode != tc.wantStatus || body.Error != tc.wantError { t.Fatalf("profile response = %d %+v", response.StatusCode, body) }
			if tc.wantStatus == 200 && (body.Data["email"] != "n@example.com" || body.Data["password"] != nil) { t.Fatalf("profile exposed unsafe data: %+v", body.Data) }
		})
	}
}

func TestRegisterAndProfileValidationResponses(t *testing.T) {
	for _, tc := range []struct { path, body, want string }{
		{"/api/auth/register", `{`, "invalid request body"},
		{"/api/auth/register", `{"name":"N","email":"n@example.com"}`, "name, email, and password are required"},
	} {
		response := authRequest(t, authApp("POST", tc.path, NewHandler(workflowsStub{}).Register), "POST", tc.path, tc.body)
		var body map[string]any
		_ = json.NewDecoder(response.Body).Decode(&body)
		response.Body.Close()
		if response.StatusCode != 400 || body["error"] != tc.want { t.Errorf("register response = %d %v", response.StatusCode, body) }
	}
	response := authRequest(t, authApp("GET", "/api/auth/me", NewHandler(workflowsStub{profileErr: errors.New("database unavailable")}).Me), "GET", "/api/auth/me", "")
	defer response.Body.Close()
	var body map[string]any
	_ = json.NewDecoder(response.Body).Decode(&body)
	if response.StatusCode != 500 || body["error"] != "failed to load user profile" { t.Fatalf("profile error response = %d %v", response.StatusCode, body) }
}

func TestRegisterMapsServiceFailures(t *testing.T) {
	for _, tc := range []struct {
		name string
		err error
		status int
		message string
	}{
		{"duplicate email", application.ErrEmailAlreadyRegistered, 409, "email already registered"},
		{"hash failure", application.ErrPasswordHashing, 500, "failed to hash password"},
		{"token failure", application.ErrTokenGeneration, 500, "failed to generate token"},
		{"unexpected failure", errors.New("database unavailable"), 500, "failed to register user"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			handler := NewHandler(workflowsStub{registerErr: tc.err})
			response := authRequest(t, authApp("POST", "/api/auth/register", handler.Register), "POST", "/api/auth/register", `{"name":"N","email":"n@example.com","password":"pw"}`)
			defer response.Body.Close()
			var body map[string]any
			if err := json.NewDecoder(response.Body).Decode(&body); err != nil { t.Fatal(err) }
			if response.StatusCode != tc.status || body["error"] != tc.message {
				t.Fatalf("register response = %d %v, want %d %q", response.StatusCode, body, tc.status, tc.message)
			}
		})
	}
}

func TestLogoutExpiresAuthCookie(t *testing.T) {
	response := authRequest(t, authApp("POST", "/api/auth/logout", NewHandler(workflowsStub{}).Logout), "POST", "/api/auth/logout", "")
	defer response.Body.Close()
	cookies := response.Cookies()
	if response.StatusCode != 200 || len(cookies) != 1 || cookies[0].Name != "jwt" || cookies[0].Value != "" || !cookies[0].HttpOnly { t.Fatalf("logout response = %d cookie=%+v", response.StatusCode, cookies) }
}
