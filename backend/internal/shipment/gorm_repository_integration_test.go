package shipment

import (
	"fmt"
	"net/url"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/narinsak-u/backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestGormRepositoryTransactionsAndConcurrentIDs(t *testing.T) {
	dsn := os.Getenv("SHIPMENT_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("set SHIPMENT_TEST_DATABASE_URL to an expendable PostgreSQL database")
	}
	base, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil { t.Fatal(err) }
	schema := fmt.Sprintf("shipment_test_%d", time.Now().UnixNano())
	if err := base.Exec("CREATE SCHEMA " + schema).Error; err != nil { t.Fatal(err) }
	t.Cleanup(func() { base.Exec("DROP SCHEMA " + schema + " CASCADE") })
	u, err := url.Parse(dsn)
	if err != nil { t.Fatal(err) }
	query := u.Query()
	query.Set("search_path", schema)
	u.RawQuery = query.Encode()
	openDB := func() *gorm.DB {
		db, err := gorm.Open(postgres.Open(u.String()), &gorm.Config{})
		if err != nil { t.Fatal(err) }
		return db
	}
	db := openDB()
	if err := db.AutoMigrate(&models.Shipment{}, &models.ShipmentEvent{}); err != nil {
		t.Fatal(err)
	}
	if err := db.Exec("ALTER TABLE shipment_events ADD CONSTRAINT reject_forced_event CHECK (status <> 'force-fail')").Error; err != nil {
		t.Fatal(err)
	}
	first, second := NewGormRepository(db), NewGormRepository(openDB())

	shipment := testShipment()
	if err := first.CreateWithEvent(&shipment, &models.ShipmentEvent{Status: "force-fail", Location: models.Location{Name: "Origin"}}); err == nil {
		t.Fatal("expected event insert failure")
	}
	var count int64
	if err := db.Model(&models.Shipment{}).Count(&count).Error; err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("shipment persisted after event failure: %d", count)
	}

	mapped := testShipment()
	mapped.Customer.Coords = models.GeoPoint{Lat: 13.75, Lng: 100.5}
	mapped.Receiver.Coords = models.GeoPoint{Lat: 18.79, Lng: 98.98}
	mapped.CurrentCoords = models.GeoPoint{Lat: 15.87, Lng: 100.99}
	if err := first.CreateWithEvent(&mapped, &models.ShipmentEvent{Status: "Label Created"}); err != nil {
		t.Fatal(err)
	}
	loaded, err := first.FindByOrderID(mapped.OrderID)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Customer.Coords != mapped.Customer.Coords || loaded.Receiver.Coords != mapped.Receiver.Coords || loaded.CurrentCoords != mapped.CurrentCoords {
		t.Fatalf("coordinate mapping did not round-trip: %+v", loaded)
	}
	duplicate := models.Shipment{OrderID: mapped.OrderID, TrackingNumber: "TH202699999", Origin: "Origin", Destination: "Destination"}
	if err := db.Create(&duplicate).Error; err == nil {
		t.Fatal("expected unique order ID constraint")
	}
	duplicate = models.Shipment{OrderID: "ORD-99999", TrackingNumber: mapped.TrackingNumber, Origin: "Origin", Destination: "Destination"}
	if err := db.Create(&duplicate).Error; err == nil {
		t.Fatal("expected unique tracking number constraint")
	}

	const concurrentCreates = 12
	var wg sync.WaitGroup
	errs := make(chan error, concurrentCreates)
	ids := make(chan [2]string, concurrentCreates)
	for i := range concurrentCreates {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			repo := first
			if i%2 != 0 { repo = second }
			shipment := testShipment()
			event := &models.ShipmentEvent{Status: "Label Created", Location: models.Location{Name: "Origin"}}
			if err := repo.CreateWithEvent(&shipment, event); err != nil { errs <- err; return }
			ids <- [2]string{shipment.OrderID, shipment.TrackingNumber}
		}(i)
	}
	wg.Wait()
	close(errs)
	for err := range errs { t.Error(err) }
	close(ids)
	orders, tracking := map[string]bool{}, map[string]bool{}
	for id := range ids {
		if orders[id[0]] || tracking[id[1]] { t.Errorf("duplicate allocated IDs: %v", id) }
		orders[id[0]], tracking[id[1]] = true, true
	}
	if len(orders) != concurrentCreates || len(tracking) != concurrentCreates {
		t.Fatalf("got %d order IDs and %d tracking IDs", len(orders), len(tracking))
	}
}

func testShipment() models.Shipment {
	return models.Shipment{
		Origin: "Origin", Destination: "Destination", Status: "pending",
		Customer: models.ContactInfo{Name: "Sender"}, Receiver: models.ContactInfo{Name: "Recipient"},
	}
}
