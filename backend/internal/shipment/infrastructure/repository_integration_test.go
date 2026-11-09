package infrastructure

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/narinsak-u/backend/internal/shipment/application"
	"github.com/narinsak-u/backend/internal/shipment/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestGormRepositoryTransactionsAndConcurrentIDs(t *testing.T){
	dsn:=os.Getenv("SHIPMENT_TEST_DATABASE_URL");if dsn==""{t.Skip("set SHIPMENT_TEST_DATABASE_URL to an expendable PostgreSQL database")}
	base,err:=gorm.Open(postgres.Open(dsn),&gorm.Config{});if err!=nil{t.Fatal(err)}
	schema:=fmt.Sprintf("shipment_test_%d",time.Now().UnixNano());if err:=base.Exec("CREATE SCHEMA "+schema).Error;err!=nil{t.Fatal(err)};t.Cleanup(func(){base.Exec("DROP SCHEMA "+schema+" CASCADE")})
	u,err:=url.Parse(dsn);if err!=nil{t.Fatal(err)};query:=u.Query();query.Set("search_path",schema);u.RawQuery=query.Encode();openDB:=func()*gorm.DB{db,err:=gorm.Open(postgres.Open(u.String()),&gorm.Config{});if err!=nil{t.Fatal(err)};return db}
	db:=openDB();if err:=db.AutoMigrate(&shipmentRecord{},&eventRecord{});err!=nil{t.Fatal(err)};if err:=db.Exec("ALTER TABLE shipment_events ADD CONSTRAINT reject_forced_event CHECK (status <> 'force-fail')").Error;err!=nil{t.Fatal(err)}
	first,second:=NewGormRepository(db),NewGormRepository(openDB());ctx:=context.Background();shipment:=testShipment()
	if err:=first.CreateWithEvent(ctx,&shipment,&domain.ShipmentEvent{Status:"force-fail",Location:domain.Location{Name:"Origin"}});err==nil{t.Fatal("expected event insert failure")};var count int64;if err:=db.Model(&shipmentRecord{}).Count(&count).Error;err!=nil{t.Fatal(err)};if count!=0{t.Fatalf("shipment persisted after event failure: %d",count)}
	mapped:=testShipment();mapped.Customer.Coords=domain.GeoPoint{Lat:13.75,Lng:100.5};mapped.Receiver.Coords=domain.GeoPoint{Lat:18.79,Lng:98.98};mapped.CurrentCoords=domain.GeoPoint{Lat:15.87,Lng:100.99};if err:=first.CreateWithEvent(ctx,&mapped,&domain.ShipmentEvent{Status:"Label Created"});err!=nil{t.Fatal(err)};loaded,err:=first.FindByOrderID(ctx,mapped.OrderID);if err!=nil{t.Fatal(err)};if loaded.Customer.Coords!=mapped.Customer.Coords||loaded.Receiver.Coords!=mapped.Receiver.Coords||loaded.CurrentCoords!=mapped.CurrentCoords{t.Fatalf("coordinate mapping did not round-trip: %+v",loaded)}
	loaded.Status=domain.StatusInTransit;if err:=first.SaveWithEvent(ctx,loaded,&domain.ShipmentEvent{Status:"force-fail"});err==nil{t.Fatal("expected status event insert failure")};loaded,err=first.FindByOrderID(ctx,mapped.OrderID);if err!=nil{t.Fatal(err)};if loaded.Status!=domain.StatusPending{t.Fatalf("status persisted after event failure: %q",loaded.Status)}
	if err:=db.Exec("CREATE TABLE shipment_delete_guard (shipment_id bigint REFERENCES shipments(id))").Error;err!=nil{t.Fatal(err)};if err:=db.Exec("INSERT INTO shipment_delete_guard (shipment_id) VALUES (?)",mapped.ID).Error;err!=nil{t.Fatal(err)};if err:=first.DeleteWithEvents(ctx,loaded);err==nil{t.Fatal("expected shipment delete failure")};if err:=db.Model(&eventRecord{}).Where("shipment_id = ?",loaded.ID).Count(&count).Error;err!=nil{t.Fatal(err)};if count!=1{t.Fatalf("events deleted after shipment delete failure: %d",count)}
	duplicate:=shipmentRecord{OrderID:mapped.OrderID,TrackingNumber:"TH202699999",Origin:"Origin",Destination:"Destination"};if err:=db.Create(&duplicate).Error;err==nil{t.Fatal("expected unique order ID constraint")};duplicate=shipmentRecord{OrderID:"ORD-99999",TrackingNumber:mapped.TrackingNumber,Origin:"Origin",Destination:"Destination"};if err:=db.Create(&duplicate).Error;err==nil{t.Fatal("expected unique tracking number constraint")}
	const concurrentCreates=12;var wg sync.WaitGroup;errs:=make(chan error,concurrentCreates);ids:=make(chan [2]string,concurrentCreates);for i:=range concurrentCreates{wg.Add(1);go func(i int){defer wg.Done();repo:=first;if i%2!=0{repo=second};s:=testShipment();if err:=repo.CreateWithEvent(ctx,&s,&domain.ShipmentEvent{Status:"Label Created",Location:domain.Location{Name:"Origin"}});err!=nil{errs<-err;return};ids<-[2]string{s.OrderID,s.TrackingNumber}}(i)};wg.Wait();close(errs);for err:=range errs{t.Error(err)};close(ids);orders,tracking:=map[string]bool{},map[string]bool{};for id:=range ids{if orders[id[0]]||tracking[id[1]]{t.Errorf("duplicate allocated IDs: %v",id)};orders[id[0]],tracking[id[1]]=true,true};if len(orders)!=concurrentCreates||len(tracking)!=concurrentCreates{t.Fatalf("got %d order IDs and %d tracking IDs",len(orders),len(tracking))}
}
func testShipment()domain.Shipment{return domain.Shipment{Origin:"Origin",Destination:"Destination",Status:domain.StatusPending,Customer:domain.Contact{Name:"Sender"},Receiver:domain.Contact{Name:"Recipient"}}}
var _ application.Repository = (*GormRepository)(nil)
