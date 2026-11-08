package seed

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/narinsak-u/backend/internal/auth"
	"github.com/narinsak-u/backend/internal/config"
	"github.com/narinsak-u/backend/internal/hub"
	"github.com/narinsak-u/backend/internal/middleware"
	"github.com/narinsak-u/backend/internal/models"
	"github.com/narinsak-u/backend/pkg/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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
	if err := db.AutoMigrate(&models.User{}, &models.Shipment{}, &models.ShipmentEvent{}, &models.Hub{}, &bootstrapState{}); err != nil {
		t.Fatal(err)
	}

	legacyHash, err := utils.HashPassword("password")
	if err != nil {
		t.Fatal(err)
	}
	legacy := models.User{Name: "Legacy", Email: "legacy@example.test", Password: legacyHash, Role: "admin"}
	otherLegacy := models.User{Name: "Other Legacy", Email: "other-legacy@example.test", Password: legacyHash, Role: "admin"}
	if err := db.Create(&legacy).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&otherLegacy).Error; err != nil {
		t.Fatal(err)
	}
	shipment := models.Shipment{OrderID: "ORDER-PRESERVED", TrackingNumber: "TRACK-PRESERVED", Origin: "A", Destination: "B", Status: "pending"}
	if err := db.Create(&shipment).Error; err != nil {
		t.Fatal(err)
	}
	preservedHub := models.Hub{ID: "HUB-PRESERVED", Name: "Preserved", CarrierID: "carrier", Address: "address", Capacity: 10, CurrentUtilization: 1, Status: "active"}
	if err := db.Create(&preservedHub).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&models.ShipmentEvent{ShipmentID: shipment.ID, Status: "pending", Location: models.Location{Name: "A"}}).Error; err != nil {
		t.Fatal(err)
	}

	oldConfig := config.App
	config.App.JWTSecret = "test-only-secret"
	config.App.JWTTTL = time.Hour
	t.Cleanup(func() { config.App = oldConfig })
	repo := auth.NewGormRepository(db)
	hubHandler := hub.NewHandler(hub.NewGormRepository(db))
	app := fiber.New()
	authHandler := auth.NewHandler(repo)
	app.Post("/api/auth/login", authHandler.Login)
	app.Post("/api/hubs", middleware.AuthRequired(), middleware.StaffOnly(repo), hubHandler.Create)
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
	status, oldToken := login("legacy@example.test", "password")
	if status != 200 || oldToken == "" {
		t.Fatalf("legacy login before cutover: status %d", status)
	}

	if err := BootstrapAdmin(db, "New Staff", "admin@example.test", "new-secret", false); err == nil {
		t.Fatal("bootstrap ran without explicit cutover opt-in")
	}
	var beforeOptIn models.User
	if err := db.First(&beforeOptIn, legacy.ID).Error; err != nil || beforeOptIn.Role != "admin" {
		t.Fatalf("disabled cutover altered legacy role: %q, err %v", beforeOptIn.Role, err)
	}
	collision := models.User{Name: "Unrelated", Email: "admin@example.test", Password: legacyHash, Role: "customer"}
	if err := db.Create(&collision).Error; err != nil {
		t.Fatal(err)
	}
	if err := BootstrapAdmin(db, "New Staff", "admin@example.test", "new-secret", true); err == nil {
		t.Fatal("bootstrap accepted an existing email")
	}
	var unchanged models.User
	if err := db.First(&unchanged, legacy.ID).Error; err != nil || unchanged.Role != "admin" {
		t.Fatalf("collision altered legacy role: role %q, err %v", unchanged.Role, err)
	}
	if err := db.Delete(&collision).Error; err != nil {
		t.Fatal(err)
	}

	if err := BootstrapAdmin(db, "New Staff", "admin@example.test", "new-secret", true); err != nil {
		t.Fatal(err)
	}
	if got := createHub(oldToken); got != 403 {
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
	}{{&models.Shipment{}, "shipment", 1}, {&models.Hub{}, "hub", 2}, {&models.ShipmentEvent{}, "event", 1}} {
		var count int64
		if err := db.Model(record.model).Count(&count).Error; err != nil || count != record.want {
			t.Fatalf("preserved %s count = %d, want %d, err %v", record.label, count, record.want, err)
		}
	}
	var legacyAfter models.User
	if err := db.First(&legacyAfter, legacy.ID).Error; err != nil || legacyAfter.Role != "customer" {
		t.Fatalf("legacy role after cutover = %q, err %v", legacyAfter.Role, err)
	}
	var otherLegacyAfter models.User
	if err := db.First(&otherLegacyAfter, otherLegacy.ID).Error; err != nil || otherLegacyAfter.Role != "customer" {
		t.Fatalf("second legacy role after cutover = %q, err %v", otherLegacyAfter.Role, err)
	}
}
