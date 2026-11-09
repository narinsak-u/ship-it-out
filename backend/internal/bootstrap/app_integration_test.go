package bootstrap

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	authinfra "github.com/narinsak-u/backend/internal/auth/infrastructure"
	hubinfra "github.com/narinsak-u/backend/internal/hub/infrastructure"
	"github.com/narinsak-u/backend/internal/platform/config"
	shipmentinfra "github.com/narinsak-u/backend/internal/shipment/infrastructure"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestNewAppComposesPublicAndProtectedRoutes(t *testing.T) {
	dsn := os.Getenv("ADMIN_BOOTSTRAP_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("set ADMIN_BOOTSTRAP_TEST_DATABASE_URL to an expendable PostgreSQL database")
	}
	base, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	schema := fmt.Sprintf("route_composition_test_%d", time.Now().UnixNano())
	if err := base.Exec("CREATE SCHEMA " + schema).Error; err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { base.Exec("DROP SCHEMA " + schema + " CASCADE") })
	databaseURL, err := url.Parse(dsn)
	if err != nil {
		t.Fatal(err)
	}
	query := databaseURL.Query()
	query.Set("search_path", schema)
	databaseURL.RawQuery = query.Encode()
	db, err := gorm.Open(postgres.Open(databaseURL.String()), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	models := append([]any{&authinfra.UserRecord{}, &hubinfra.Record{}}, shipmentinfra.MigrationModels()...)
	if err := db.AutoMigrate(models...); err != nil {
		t.Fatal(err)
	}

	oldConfig := config.App
	config.App.JWTSecret = "route-composition-test-secret"
	config.App.JWTTTL = time.Hour
	t.Cleanup(func() { config.App = oldConfig })
	app := NewApp(db)

	for _, route := range []struct{ method, path string }{
		{"GET", "/api/health"},
		{"GET", "/api/hubs"},
		{"GET", "/api/shipments"},
		{"GET", "/api/analytics/overview"},
		{"GET", "/api/analytics/timeseries"},
	} {
		response := apiRequest(t, app, route.method, route.path, "", nil)
		var body struct { Success bool `json:"success"` }
		_ = json.NewDecoder(response.Body).Decode(&body)
		response.Body.Close()
		if response.StatusCode != 200 || !body.Success {
			t.Errorf("%s %s = %d, success=%t", route.method, route.path, response.StatusCode, body.Success)
		}
	}
	for _, route := range []struct{ method, path, want string }{
		{"GET", "/api/hubs/HUB-MISSING", "hub not found"},
		{"GET", "/api/shipments/ORDER-MISSING", "shipment not found"},
		{"GET", "/api/track/TRACK-MISSING", "shipment not found"},
	} {
		response := apiRequest(t, app, route.method, route.path, "", nil)
		var body struct { Error string `json:"error"` }
		_ = json.NewDecoder(response.Body).Decode(&body)
		response.Body.Close()
		if response.StatusCode != 404 || body.Error != route.want {
			t.Errorf("%s %s = %d %q, want 404 %q", route.method, route.path, response.StatusCode, body.Error, route.want)
		}
	}

	for _, route := range []struct{ method, path string }{
		{"POST", "/api/hubs"}, {"PUT", "/api/hubs/HUB-MISSING"}, {"DELETE", "/api/hubs/HUB-MISSING"},
		{"POST", "/api/shipments"}, {"PATCH", "/api/shipments/ORDER-MISSING/status"},
		{"PUT", "/api/shipments/ORDER-MISSING"}, {"DELETE", "/api/shipments/ORDER-MISSING"},
	} {
		response := apiRequest(t, app, route.method, route.path, `{}`, nil)
		response.Body.Close()
		if response.StatusCode != 401 {
			t.Errorf("unauthenticated %s %s = %d, want 401", route.method, route.path, response.StatusCode)
		}
	}

	register := apiRequest(t, app, "POST", "/api/auth/register", `{"name":"Customer","email":"route@example.test","password":"secret"}`, nil)
	cookies := register.Cookies()
	register.Body.Close()
	if register.StatusCode != 200 || len(cookies) != 1 || cookies[0].Name != "jwt" {
		t.Fatalf("registration response = %d cookies=%+v", register.StatusCode, cookies)
	}
	login := apiRequest(t, app, "POST", "/api/auth/login", `{"email":"route@example.test","password":"secret"}`, nil)
	loginCookies := login.Cookies()
	var loginBody struct {
		Success bool `json:"success"`
		Data struct { Token string `json:"token"` } `json:"data"`
	}
	_ = json.NewDecoder(login.Body).Decode(&loginBody)
	login.Body.Close()
	if login.StatusCode != 200 || !loginBody.Success || loginBody.Data.Token == "" || len(loginCookies) != 1 {
		t.Fatalf("login response = %d %+v cookies=%+v", login.StatusCode, loginBody, loginCookies)
	}
	logout := apiRequest(t, app, "POST", "/api/auth/logout", "", loginCookies[0])
	logoutCookies := logout.Cookies()
	logout.Body.Close()
	if logout.StatusCode != 200 || len(logoutCookies) != 1 || logoutCookies[0].Name != "jwt" || logoutCookies[0].Value != "" {
		t.Fatalf("logout response = %d cookies=%+v", logout.StatusCode, logoutCookies)
	}
	profile := apiRequest(t, app, "GET", "/api/auth/me", "", cookies[0])
	profile.Body.Close()
	if profile.StatusCode != 200 {
		t.Fatalf("authenticated profile status = %d, want 200", profile.StatusCode)
	}
	for _, route := range []struct{ method, path string }{
		{"POST", "/api/hubs"}, {"PUT", "/api/hubs/HUB-MISSING"}, {"DELETE", "/api/hubs/HUB-MISSING"},
		{"POST", "/api/shipments"}, {"PATCH", "/api/shipments/ORDER-MISSING/status"},
		{"PUT", "/api/shipments/ORDER-MISSING"}, {"DELETE", "/api/shipments/ORDER-MISSING"},
	} {
		response := apiRequest(t, app, route.method, route.path, `{}`, cookies[0])
		response.Body.Close()
		if response.StatusCode != 403 {
			t.Errorf("customer %s %s = %d, want 403", route.method, route.path, response.StatusCode)
		}
	}
}

func apiRequest(t *testing.T, app interface{ Test(*http.Request, ...int) (*http.Response, error) }, method, path, body string, cookie *http.Cookie) *http.Response {
	t.Helper()
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" { request.Header.Set("Content-Type", "application/json") }
	if cookie != nil { request.AddCookie(cookie) }
	response, err := app.Test(request)
	if err != nil { t.Fatal(err) }
	return response
}
