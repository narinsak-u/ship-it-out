package presentation

import (
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/narinsak-u/backend/internal/shipment/domain"
	trackingapp "github.com/narinsak-u/backend/internal/tracking/application"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type trackingReader struct { shipment *domain.Shipment; events []domain.ShipmentEvent; shipmentErr, eventsErr error }
func (r trackingReader) FindByTrackingNumber(context.Context, string) (*domain.Shipment, error) { return r.shipment, r.shipmentErr }
func (r trackingReader) FindEventsByShipmentID(context.Context, uint) ([]domain.ShipmentEvent, error) { return r.events, r.eventsErr }

func TestTrackJSONContract(t *testing.T) {
	created := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	s := &domain.Shipment{ID: 7, OrderID: "ORD-7", TrackingNumber: "TH-7", Customer: domain.Contact{Name:"A", Zipcode:"10000", SubDistrict:"Sub", District:"Dist", Province:"Prov", Coords:domain.GeoPoint{Lat:13.5,Lng:100.5}}, Receiver:domain.Contact{Name:"B",Coords:domain.GeoPoint{Lat:14,Lng:101}}, CurrentCoords:domain.GeoPoint{Lat:13,Lng:100}, Status:domain.StatusInTransit, CreatedAt:created, EstimatedDelivery:created}
	events := []domain.ShipmentEvent{{ID:11, ShipmentID:7, Status:"in_transit", Location:domain.Location{Name:"Hub",Lat:13,Lng:100}, Description:"Moving", CreatedAt:created}, {ID:12, ShipmentID:7, Status:"delivered", Location:domain.Location{Name:"End",Lat:14,Lng:101}, CreatedAt:created.Add(time.Hour)}}
	app := fiber.New()
	app.Get("/api/track/:trackingNumber", NewHandler(trackingapp.NewService(trackingReader{shipment:s, events:events})).Track)
	resp, err := app.Test(httptest.NewRequest("GET", "/api/track/TH-7", nil))
	require.NoError(t, err); defer resp.Body.Close(); assert.Equal(t, 200, resp.StatusCode)
	var got map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&got))
	want := `{"success":true,"data":{"shipment":{"id":"ORD-7","trackingNumber":"TH-7","customer":{"name":"A","zipcode":"10000","subDistrict":"Sub","district":"Dist","province":"Prov","coords":{"lat":13.5,"lng":100.5}},"receiver":{"name":"B","zipcode":"","subDistrict":"","district":"","province":"","coords":{"lat":14,"lng":101}},"currentCoords":{"lat":13,"lng":100},"origin":"","destination":"","status":"in_transit","carrier":"","weight":0,"items":0,"estimatedDelivery":"2026-01-02T03:04:05Z","createdAt":"2026-01-02T03:04:05Z","progress":0},"events":[{"id":11,"shipmentId":7,"status":"in_transit","location":{"name":"Hub","lat":13,"lng":100},"description":"Moving","timestamp":"2026-01-02T03:04:05Z"},{"id":12,"shipmentId":7,"status":"delivered","location":{"name":"End","lat":14,"lng":101},"timestamp":"2026-01-02T04:04:05Z"}]}}`
	var expected map[string]any
	require.NoError(t, json.Unmarshal([]byte(want), &expected))
	assert.Equal(t, expected, got)
}

func TestTrackErrors(t *testing.T) {
	for _, tc := range []struct{name string; reader trackingReader; status int}{
		{name:"not found", reader:trackingReader{shipmentErr:domain.ErrNotFound}, status:404},
		{name:"repository failure", reader:trackingReader{shipmentErr:errors.New("db")}, status:500},
		{name:"event failure", reader:trackingReader{shipment:&domain.Shipment{ID:1}, eventsErr:errors.New("db")}, status:500},
	} { t.Run(tc.name, func(t *testing.T) {
		app := fiber.New()
		app.Get("/api/track/:trackingNumber", NewHandler(trackingapp.NewService(tc.reader)).Track)
		resp, err := app.Test(httptest.NewRequest("GET", "/api/track/TH-7", nil))
		require.NoError(t, err); defer resp.Body.Close(); assert.Equal(t, tc.status, resp.StatusCode)
	}) }
}
