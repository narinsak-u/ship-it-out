package application

import (
	"context"
	"errors"
	"testing"
	"time"

	hubdomain "github.com/narinsak-u/backend/internal/hub/domain"
	"github.com/narinsak-u/backend/internal/shipment/domain"
)

type repositoryFake struct { created *domain.Shipment; createdEvent *domain.ShipmentEvent; saved *domain.Shipment; savedEvent *domain.ShipmentEvent; deleted bool; found *domain.Shipment; err error }
func (r *repositoryFake) List(context.Context,domain.Filter)([]domain.Shipment,int64,error){return nil,0,r.err}
func (r *repositoryFake) FindByOrderID(context.Context,string)(*domain.Shipment,error){if r.err!=nil{return nil,r.err};return r.found,nil}
func (r *repositoryFake) CreateWithEvent(_ context.Context,s *domain.Shipment,e *domain.ShipmentEvent)error{r.created=s;r.createdEvent=e;return r.err}
func (r *repositoryFake) SaveWithEvent(_ context.Context,s *domain.Shipment,e *domain.ShipmentEvent)error{r.saved=s;r.savedEvent=e;return r.err}
func (r *repositoryFake) DeleteWithEvents(context.Context,*domain.Shipment)error{r.deleted=true;return r.err}
type hubFake struct { hub *hubdomain.Hub; err error }
func(h hubFake)FindByID(context.Context,string)(*hubdomain.Hub,error){return h.hub,h.err}

func TestCreatePersistsShipmentAndInitialEvent(t *testing.T){now:=time.Date(2026,1,2,3,4,5,0,time.UTC);repo:=&repositoryFake{};service:=NewService(repo,nil);service.now=func()time.Time{return now};sender:=domain.Contact{Name:"Sender",Province:"Bangkok",Zipcode:"10110",Coords:domain.GeoPoint{Lat:13.7,Lng:100.5}};receiver:=domain.Contact{Name:"Receiver",Province:"Chonburi",Zipcode:"20110"};got,err:=service.Create(context.Background(),CreateInput{Customer:sender,Receiver:receiver,Carrier:"Express",Weight:2,Items:1});if err!=nil{t.Fatal(err)};if got.Status!="pending"||got.Progress!=0||got.EstimatedDelivery!=now.Add(72*time.Hour){t.Fatalf("unexpected shipment: %+v",got)};if repo.created!=got||repo.createdEvent==nil||repo.createdEvent.Status!="Label Created"||repo.createdEvent.Description!="Awaiting pickup."||repo.createdEvent.Location.Lat!=13.7{t.Fatalf("shipment and initial event not persisted together: %+v %+v",repo.created,repo.createdEvent)}}
func TestCreateRejectsMissingCustomer(t *testing.T){_,err:=NewService(&repositoryFake{},nil).Create(context.Background(),CreateInput{});if err==nil||err.Error()!="invalid shipment input: customer name, province, and zipcode are required"{t.Fatalf("unexpected validation error: %v",err)}}
func TestUpdateStatusPersistsHubAndMappedEvent(t *testing.T){shipment:=&domain.Shipment{ID:9,OrderID:"ORD-10246",Status:domain.StatusPending,Origin:"Origin",Destination:"Destination",Customer:domain.Contact{Name:"Sender",Province:"Bangkok",Zipcode:"10110"},CurrentCoords:domain.GeoPoint{Lat:1,Lng:2}};repo:=&repositoryFake{found:shipment};hub:=&hubdomain.Hub{ID:"HUB-1",Name:"Central",Address:"Road",Coords:hubdomain.GeoPoint{Lat:14,Lng:101}};got,err:=NewService(repo,hubFake{hub:hub}).UpdateStatus(context.Background(),shipment.OrderID,domain.StatusInTransit,hub.ID);if err!=nil{t.Fatal(err)};if got.HubID!=hub.ID||got.CurrentCoords!=(domain.GeoPoint{Lat:14,Lng:101})||repo.savedEvent==nil||repo.savedEvent.Status!="In Transit"||repo.savedEvent.Location.Name!="Central, Road"{t.Fatalf("unexpected update/event: %+v %+v",got,repo.savedEvent)}}
func TestUpdateStatusRejectsMissingHub(t *testing.T){repo:=&repositoryFake{found:&domain.Shipment{OrderID:"ORD-1"}};_,err:=NewService(repo,hubFake{err:domain.ErrNotFound}).UpdateStatus(context.Background(),"ORD-1",domain.StatusInTransit,"HUB-X");if !errors.Is(err,ErrInvalidHub){t.Fatalf("expected invalid hub, got %v",err)}}
