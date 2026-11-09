package infrastructure

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/narinsak-u/backend/internal/hub/application"
	"github.com/narinsak-u/backend/internal/hub/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Hub = domain.Hub
type GeoPoint = domain.GeoPoint
var ErrNotFound = application.ErrNotFound

func TestGormRepositoryCRUD(t *testing.T) {
	dsn := os.Getenv("HUB_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("set HUB_TEST_DATABASE_URL to an expendable PostgreSQL database")
	}
	base, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	schema := fmt.Sprintf("hub_test_%d", time.Now().UnixNano())
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
	if err := db.AutoMigrate(&Record{}); err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	repo := NewGormRepository(db)
	first := Hub{ID: "HUB-004", Name: "Bangkok Hub", CarrierID: "THUN", Address: "Bangkok", Coords: GeoPoint{Lat: 13.7563, Lng: 100.5018}, Capacity: 5000, CurrentUtilization: 3400, Status: "active"}
	if err := repo.Create(ctx, &first); err != nil {
		t.Fatal(err)
	}
	generated := Hub{Name: "Chonburi Hub", CarrierID: "THUN", Address: "Chonburi", Coords: GeoPoint{Lat: 13.3611, Lng: 100.9847}, Capacity: 3000, CurrentUtilization: 1200, Status: "active"}
	if err := repo.Create(ctx, &generated); err != nil {
		t.Fatal(err)
	}
	if generated.ID != "HUB-005" {
		t.Fatalf("generated ID = %q, want HUB-005", generated.ID)
	}

	loaded, err := repo.FindByID(ctx, first.ID)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Coords != first.Coords {
		t.Fatalf("coordinates = %+v, want %+v", loaded.Coords, first.Coords)
	}
	loaded.Status = "maintenance"
	if err := repo.Save(ctx, loaded); err != nil {
		t.Fatal(err)
	}
	all, err := repo.FindAll(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 2 {
		t.Fatalf("listed %d hubs, want 2", len(all))
	}
	foundUpdated := false
	for _, listed := range all {
		if listed.ID == first.ID && listed.Status == "maintenance" {
			foundUpdated = true
		}
	}
	if !foundUpdated {
		t.Fatalf("updated hub missing from list: %+v", all)
	}
	if err := repo.Delete(ctx, generated.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.FindByID(ctx, generated.ID); err != ErrNotFound {
		t.Fatalf("lookup after delete error = %v, want %v", err, ErrNotFound)
	}
}
