package seed

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"context"
	"fmt"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"errors"
	authapplication "github.com/narinsak-u/backend/internal/auth/application"
	authinfra "github.com/narinsak-u/backend/internal/auth/infrastructure"
	authpresentation "github.com/narinsak-u/backend/internal/auth/presentation"
	"github.com/narinsak-u/backend/internal/platform/config"
	hubapplication "github.com/narinsak-u/backend/internal/hub/application"
	hubinfra "github.com/narinsak-u/backend/internal/hub/infrastructure"
	hubpresentation "github.com/narinsak-u/backend/internal/hub/presentation"
	hubdomain "github.com/narinsak-u/backend/internal/hub/domain"
	shipmentapplication "github.com/narinsak-u/backend/internal/shipment/application"
	shipmentdomain "github.com/narinsak-u/backend/internal/shipment/domain"
	shipmentinfra "github.com/narinsak-u/backend/internal/shipment/infrastructure"
	shipmentpresentation "github.com/narinsak-u/backend/internal/shipment/presentation"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type testShipmentRecord struct {
	ID             uint `gorm:"primaryKey"`
	OrderID        string
	TrackingNumber string
	Origin         string
	Destination    string
	Status         string
}

func (testShipmentRecord) TableName() string { return "shipments" }

type testShipmentEventRecord struct {
	ID         uint `gorm:"primaryKey"`
	ShipmentID uint
	Status     string
	LocationName string
}

func (testShipmentEventRecord) TableName() string { return "shipment_events" }
func TestAdminCutoverHTTPAndPersistence(t *testing.T) {

	dsn := os.Getenv("ADMIN_BOOTSTRAP_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("set ADMIN_BOOTSTRAP_TEST_DATABASE_URL to an expendable PostgreSQL database")
	}
	base, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	schema := fmt.Sprintf("admin_cutover_test_%d", time.Now().UnixNano())
	if err := base.Exec("CREATE SCHEMA " + schema).Error; err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { base.Exec("DROP SCHEMA " + schema + " CASCADE") })
	u, err := url.Parse(dsn)
	if err != nil {
		t.Fatal(err)
	}
	query := u.Query()
	query.Set("search_path", schema)
	u.RawQuery = query.Encode()
	db, err := gorm.Open(postgres.Open(u.String()), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	legacyHash, err := authinfra.HashPassword("password")
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Exec(`CREATE TABLE users (
		id BIGSERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		role TEXT NOT NULL DEFAULT 'customer',
		created_at TIMESTAMPTZ
	)`).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Exec("INSERT INTO users (name, email, password, role) VALUES (?, ?, ?, ?)", "Existing", "existing@example.test", legacyHash, "customer").Error; err != nil {
		t.Fatal(err)
	}
	migrationModels := append([]any{&authinfra.UserRecord{}, &hubinfra.Record{}, &bootstrapState{}, &demoBootstrapState{}}, shipmentinfra.MigrationModels()...)
	if err := db.AutoMigrate(migrationModels...); err != nil {
		t.Fatal(err)
	}
	existing, err := authinfra.NewGormRepository(db).FindByEmail("existing@example.test")
	if err != nil || existing.Role != "customer" || !authinfra.CheckPassword("password", existing.Password) {
		t.Fatalf("pre-existing users row = %+v, err %v", existing, err)
	}

	legacy := authinfra.UserRecord{Name: "Legacy", Email: "legacy@example.test", Password: legacyHash, Role: "admin"}
	otherLegacy := authinfra.UserRecord{Name: "Other Legacy", Email: "other-legacy@example.test", Password: legacyHash, Role: "admin"}
	if err := db.Create(&legacy).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&otherLegacy).Error; err != nil {
		t.Fatal(err)
	}
	preservedShipment := testShipmentRecord{OrderID: "ORDER-PRESERVED", TrackingNumber: "TRACK-PRESERVED", Origin: "A", Destination: "B", Status: "pending"}
	if err := db.Create(&preservedShipment).Error; err != nil {
		t.Fatal(err)
	}
	preservedHub := hubinfra.Record{ID: "HUB-PRESERVED", Name: "Preserved", CarrierID: "carrier", Address: "address", Capacity: 10, CurrentUtilization: 1, Status: "active"}
	if err := db.Create(&preservedHub).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&testShipmentEventRecord{ShipmentID: preservedShipment.ID, Status: "pending", LocationName: "A"}).Error; err != nil {
		t.Fatal(err)
	}

	oldConfig := config.App
	config.App.JWTSecret = "test-only-secret"
	config.App.JWTTTL = time.Hour
	t.Cleanup(func() { config.App = oldConfig })
	repo := authinfra.NewGormRepository(db)
	authService := authapplication.NewService(repo, authinfra.PasswordAdapter{}, authinfra.NewJWTAdapter(config.App.JWTSecret, config.App.JWTTTL))
	hubHandler := hubpresentation.NewHandler(hubapplication.NewService(hubinfra.NewGormRepository(db)))
	app := fiber.New()
	authHandler := authpresentation.NewHandler(authService)
	app.Post("/api/auth/register", authHandler.Register)
	app.Post("/api/auth/login", authHandler.Login)
	app.Get("/api/auth/me", authpresentation.AuthRequired(), authHandler.Me)
	app.Post("/api/auth/logout", authHandler.Logout)
	app.Post("/api/hubs", authpresentation.AuthRequired(), authpresentation.StaffOnly(repo), hubHandler.Create)
	app.Put("/api/hubs/:id", authpresentation.AuthRequired(), authpresentation.StaffOnly(repo), hubHandler.Update)
	app.Delete("/api/hubs/:id", authpresentation.AuthRequired(), authpresentation.StaffOnly(repo), hubHandler.Delete)
	shipmentRepo := shipmentinfra.NewGormRepository(db)
	shipmentHandler := shipmentpresentation.NewHandler(shipmentapplication.NewService(shipmentRepo, seedShipmentHubReader{service: hubapplication.NewService(hubinfra.NewGormRepository(db))}))
	app.Post("/api/shipments", authpresentation.AuthRequired(), authpresentation.StaffOnly(repo), shipmentHandler.Create)
	app.Patch("/api/shipments/:orderId/status", authpresentation.AuthRequired(), authpresentation.StaffOnly(repo), shipmentHandler.UpdateStatus)
	app.Put("/api/shipments/:orderId", authpresentation.AuthRequired(), authpresentation.StaffOnly(repo), shipmentHandler.Update)
	app.Delete("/api/shipments/:orderId", authpresentation.AuthRequired(), authpresentation.StaffOnly(repo), shipmentHandler.Delete)
	register := httptest.NewRequest("POST", "/api/auth/register", strings.NewReader(`{"name":"Elevated","email":"elevated@example.test","password":"password","role":"admin"}`))
	register.Header.Set("Content-Type", "application/json")
	registerResponse, err := app.Test(register)
	if err != nil {
		t.Fatal(err)
	}
	registerResponse.Body.Close()
	if registerResponse.StatusCode != 200 {
		t.Fatalf("public signup status = %d, want 200", registerResponse.StatusCode)
	}
	var registered authinfra.UserRecord
	if err := db.Where("email = ?", "elevated@example.test").First(&registered).Error; err != nil || registered.Role != "customer" {
		t.Fatalf("signup role = %q, err %v", registered.Role, err)
	}
	login := func(email, password string) (int, string) {
		request := httptest.NewRequest("POST", "/api/auth/login", strings.NewReader(fmt.Sprintf(`{"email":%q,"password":%q}`, email, password)))
		request.Header.Set("Content-Type", "application/json")
		response, err := app.Test(request)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()
		var body struct {
			Data struct {
				Token string `json:"token"`
			} `json:"data"`
		}
		_ = json.NewDecoder(response.Body).Decode(&body)
		return response.StatusCode, body.Data.Token
	}
	createHub := func(token string) int {
		request := httptest.NewRequest("POST", "/api/hubs", strings.NewReader(`{"id":"HUB-WRITE","name":"Write test","carrierId":"carrier","address":"test","capacity":10,"currentUtilization":1,"coords":{"lat":1,"lng":2}}`))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", "Bearer "+token)
		response, err := app.Test(request)
		if err != nil {
			t.Fatal(err)
		}
		response.Body.Close()
		return response.StatusCode
	}
	writeStatus := func(token, method, path string) int {
		request := httptest.NewRequest(method, path, strings.NewReader(`{}`))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", "Bearer "+token)
		response, err := app.Test(request)
		if err != nil {
			t.Fatal(err)
		}
		response.Body.Close()
		return response.StatusCode
	}
	status, oldToken := login("legacy@example.test", "password")
	if status != 200 || oldToken == "" {
		t.Fatalf("legacy login before cutover: status %d", status)
	}

	if err := BootstrapAdmin(db, "New Staff", "admin@example.test", "new-secret", false); err == nil {
		t.Fatal("bootstrap ran without explicit cutover opt-in")
	}
	var beforeOptIn authinfra.UserRecord
	if err := db.First(&beforeOptIn, legacy.ID).Error; err != nil || beforeOptIn.Role != "admin" {
		t.Fatalf("disabled cutover altered legacy role: %q, err %v", beforeOptIn.Role, err)
	}
	collision := authinfra.UserRecord{Name: "Unrelated", Email: "admin@example.test", Password: legacyHash, Role: "customer"}
	if err := db.Create(&collision).Error; err != nil {
		t.Fatal(err)
	}
	if err := BootstrapAdmin(db, "New Staff", "admin@example.test", "new-secret", true); err == nil {
		t.Fatal("bootstrap accepted an existing email")
	}
	var unchanged authinfra.UserRecord
	if err := db.First(&unchanged, legacy.ID).Error; err != nil || unchanged.Role != "admin" {
		t.Fatalf("collision altered legacy role: role %q, err %v", unchanged.Role, err)
	}
	if err := db.Delete(&collision).Error; err != nil {
		t.Fatal(err)
	}

	if err := BootstrapAdmin(db, "New Staff", "admin@example.test", "new-secret", true); err != nil {
		t.Fatal(err)
	}
	demoCollision := authinfra.UserRecord{Name: "Unrelated Demo Email Owner", Email: demoUserEmail, Password: legacyHash, Role: "customer"}
	if err := db.Create(&demoCollision).Error; err != nil {
		t.Fatal(err)
	}
	if err := BootstrapDemoUser(db); err == nil {
		t.Fatal("demo bootstrap accepted an existing email")
	}
	var collisionUnchanged authinfra.UserRecord
	if err := db.First(&collisionUnchanged, demoCollision.ID).Error; err != nil || collisionUnchanged.Name != demoCollision.Name || collisionUnchanged.Role != "customer" || collisionUnchanged.Password != legacyHash {
		t.Fatalf("demo collision account changed: %+v, err %v", collisionUnchanged, err)
	}
	if err := db.Delete(&demoCollision).Error; err != nil {
		t.Fatal(err)
	}
	if err := BootstrapDemoUser(db); err != nil {
		t.Fatal(err)
	}
	if err := BootstrapDemoUser(db); err != nil {
		t.Fatal(err)
	}
	status, demoToken := login(demoUserEmail, demoUserPassword)
	if status != 200 || demoToken == "" {
		t.Fatalf("seeded demo login: status %d", status)
	}
	cookieLogin := httptest.NewRequest("POST", "/api/auth/login", strings.NewReader(fmt.Sprintf(`{"email":%q,"password":%q}`, demoUserEmail, demoUserPassword)))
	cookieLogin.Header.Set("Content-Type", "application/json")
	cookieLoginResponse, err := app.Test(cookieLogin)
	if err != nil {
		t.Fatal(err)
	}
	cookies := cookieLoginResponse.Cookies()
	cookieLoginResponse.Body.Close()
	if cookieLoginResponse.StatusCode != 200 || len(cookies) != 1 || cookies[0].Name != "jwt" || !cookies[0].HttpOnly {
		t.Fatalf("demo login cookie = %#v, status %d", cookies, cookieLoginResponse.StatusCode)
	}
	refresh := httptest.NewRequest("GET", "/api/auth/me", nil)
	refresh.AddCookie(cookies[0])
	refreshResponse, err := app.Test(refresh)
	if err != nil {
		t.Fatal(err)
	}
	refreshResponse.Body.Close()
	if refreshResponse.StatusCode != 200 {
		t.Fatalf("demo session refresh status = %d, want 200", refreshResponse.StatusCode)
	}
	logout := httptest.NewRequest("POST", "/api/auth/logout", nil)
	logout.AddCookie(cookies[0])
	logoutResponse, err := app.Test(logout)
	if err != nil {
		t.Fatal(err)
	}
	logoutCookies := logoutResponse.Cookies()
	logoutResponse.Body.Close()
	if logoutResponse.StatusCode != 200 || len(logoutCookies) != 1 || logoutCookies[0].Name != "jwt" || logoutCookies[0].Value != "" {
		t.Fatalf("demo logout cookie = %#v, status %d", logoutCookies, logoutResponse.StatusCode)
	}
	unsignedRefresh, err := app.Test(httptest.NewRequest("GET", "/api/auth/me", nil))
	if err != nil {
		t.Fatal(err)
	}
	unsignedRefresh.Body.Close()
	if unsignedRefresh.StatusCode != 401 {
		t.Fatalf("unsigned session refresh status = %d, want 401", unsignedRefresh.StatusCode)
	}
	var seededDemo authinfra.UserRecord
	if err := db.Where("email = ?", demoUserEmail).First(&seededDemo).Error; err != nil {
		t.Fatal(err)
	}
	var demoCount int64
	if err := db.Model(&authinfra.UserRecord{}).Where("email = ?", demoUserEmail).Count(&demoCount).Error; err != nil || demoCount != 1 {
		t.Fatalf("demo account count = %d, err %v", demoCount, err)
	}
	for _, write := range []struct{ method, path string }{
		{"POST", "/api/hubs"},
		{"PUT", "/api/hubs/HUB-PRESERVED"},
		{"DELETE", "/api/hubs/HUB-PRESERVED"},
		{"POST", "/api/shipments"},
		{"PATCH", "/api/shipments/ORDER-PRESERVED/status"},
		{"PUT", "/api/shipments/ORDER-PRESERVED"},
		{"DELETE", "/api/shipments/ORDER-PRESERVED"},
	} {
		if got := writeStatus(demoToken, write.method, write.path); got != 403 {
			t.Errorf("demo %s %s status = %d, want 403", write.method, write.path, got)
		}
	}
	var demoAfterRepeat authinfra.UserRecord
	if err := db.First(&demoAfterRepeat, seededDemo.ID).Error; err != nil || demoAfterRepeat.Email != demoUserEmail || demoAfterRepeat.Role != "customer" || demoAfterRepeat.Password != seededDemo.Password {
		t.Fatalf("repeat bootstrap changed demo identity or credentials: %+v, err %v", demoAfterRepeat, err)
	}
	if got := writeStatus(oldToken, "POST", "/api/hubs"); got != 403 {
		t.Fatalf("legacy session write status = %d, want 403", got)
	}
	status, newToken := login("admin@example.test", "new-secret")
	if status != 200 || newToken == "" {
		t.Fatalf("seeded admin login: status %d", status)
	}
	if got := createHub(newToken); got != 200 {
		t.Fatalf("seeded admin write status = %d, want 200", got)
	}
	if err := BootstrapAdmin(db, "Changed Name", "admin@example.test", "changed-secret", false); err != nil {
		t.Fatal(err)
	}
	if err := BootstrapAdmin(db, "Unrelated", "unrelated@example.test", "other-secret", false); err == nil {
		t.Fatal("repeat bootstrap accepted a different configured email")
	}
	status, repeatToken := login("admin@example.test", "new-secret")
	if status != 200 || repeatToken == "" {
		t.Fatalf("repeat bootstrap changed seeded credentials: status %d", status)
	}

	for _, record := range []struct {
		model any
		label string
		want  int64
	}{{"shipments", "shipment", 1}, {&hubinfra.Record{}, "hub", 2}, {"shipment_events", "event", 1}} {
		var count int64
		if record.model == "shipments" || record.model == "shipment_events" {
			table := record.model.(string)
			if err := db.Table(table).Count(&count).Error; err != nil || count != record.want {
				t.Fatalf("preserved %s count = %d, want %d, err %v", record.label, count, record.want, err)
			}
		} else if err := db.Model(record.model).Count(&count).Error; err != nil || count != record.want {
			t.Fatalf("preserved %s count = %d, want %d, err %v", record.label, count, record.want, err)
		}
	}
	var legacyAfter authinfra.UserRecord
	if err := db.First(&legacyAfter, legacy.ID).Error; err != nil || legacyAfter.Role != "customer" {
		t.Fatalf("legacy role after cutover = %q, err %v", legacyAfter.Role, err)
	}
	var otherLegacyAfter authinfra.UserRecord
	if err := db.First(&otherLegacyAfter, otherLegacy.ID).Error; err != nil || otherLegacyAfter.Role != "customer" {
		t.Fatalf("second legacy role after cutover = %q, err %v", otherLegacyAfter.Role, err)
	}
}
type seedShipmentHubReader struct { service *hubapplication.Service }
func(r seedShipmentHubReader)FindByID(ctx context.Context,id string)(*hubdomain.Hub,error){hub,err:=r.service.GetByID(ctx,id);if errors.Is(err,hubapplication.ErrNotFound){return nil,shipmentdomain.ErrNotFound};return hub,err}
