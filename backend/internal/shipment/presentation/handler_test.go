package presentation

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	hubdomain "github.com/narinsak-u/backend/internal/hub/domain"
	"github.com/narinsak-u/backend/internal/shipment/application"
	"github.com/narinsak-u/backend/internal/shipment/domain"
)

type repoFake struct {
	shipment *domain.Shipment
	event    *domain.ShipmentEvent
	rows     []domain.Shipment
	total    int64
	filter   domain.Filter
	deleted  bool
	err      error
}

func (r *repoFake) List(_ context.Context, filter domain.Filter) ([]domain.Shipment, int64, error) {
	r.filter = filter
	return r.rows, r.total, r.err
}
func (r *repoFake) FindByOrderID(context.Context, string) (*domain.Shipment, error) {
	if r.err != nil {
		return nil, r.err
	}
	if r.shipment == nil {
		return nil, application.ErrNotFound
	}
	return r.shipment, nil
}
func (r *repoFake) CreateWithEvent(_ context.Context, s *domain.Shipment, e *domain.ShipmentEvent) error {
	s.OrderID = "ORD-10246"
	s.TrackingNumber = "TH202600001"
	r.shipment = s
	r.event = e
	return r.err
}
func (r *repoFake) SaveWithEvent(_ context.Context, s *domain.Shipment, e *domain.ShipmentEvent) error {
	r.shipment = s
	r.event = e
	return r.err
}
func (r *repoFake) DeleteWithEvents(context.Context, *domain.Shipment) error {
	r.deleted = true
	return r.err
}

type testResponse struct {
	*http.Response
	Code int
}

func request(app *fiber.App, method, path, body string) *testResponse {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	response, err := app.Test(req)
	if err != nil {
		panic(err)
	}
	return &testResponse{Response: response, Code: response.StatusCode}
}

func TestCreateReturnsExistingShipmentJSONShape(t *testing.T) {
	repo := &repoFake{}
	app := fiber.New()
	app.Post("/api/shipments", NewHandler(application.NewService(repo, nil)).Create)
	response := request(app, "POST", "/api/shipments", `{"customer":{"name":"Sender","province":"Bangkok","zipcode":"10110","coords":{"lat":13.7,"lng":100.5}},"receiver":{"name":"Receiver","province":"Chonburi","zipcode":"20110"},"carrier":"Express","weight":2,"items":1}`)
	if response.Code != 200 {
		t.Fatalf("status %d", response.Code)
	}
	var result struct {
		Success bool        `json:"success"`
		Data    ShipmentDTO `json:"data"`
	}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if !result.Success || result.Data.ID != "ORD-10246" || result.Data.Customer.Name != "Sender" || result.Data.CurrentCoords.Lat != 13.7 || result.Data.TrackingNumber != "TH202600001" {
		t.Fatalf("unexpected response: %+v", result)
	}
	if repo.event == nil || repo.event.Status != "Label Created" {
		t.Fatalf("initial event absent: %+v", repo.event)
	}
}
func TestCreateReturnsValidationError(t *testing.T) {
	app := fiber.New()
	app.Post("/api/shipments", NewHandler(application.NewService(&repoFake{}, nil)).Create)
	response := request(app, "POST", "/api/shipments", `{}`)
	var result struct {
		Error string `json:"error"`
	}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if response.Code != 400 || result.Error != "customer name, province, and zipcode are required" {
		t.Fatalf("unexpected response: %d %+v", response.Code, result)
	}
}
func TestListPreservesFiltersAndPagination(t *testing.T) {
	row := domain.Shipment{OrderID: "ORD-10246", TrackingNumber: "TH202600001", Status: domain.StatusInTransit}
	repo := &repoFake{rows: []domain.Shipment{row}, total: 21}
	app := fiber.New()
	app.Get("/api/shipments", NewHandler(application.NewService(repo, nil)).List)
	response := request(app, "GET", "/api/shipments?page=2&limit=10&search=ORD&status=in_transit&exclude_status=delivered", "")
	var result struct {
		Data       []ShipmentDTO                                `json:"data"`
		Pagination struct{ Page, Limit, Total, TotalPages int } `json:"pagination"`
	}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if response.Code != 200 || repo.filter != (domain.Filter{Page: 2, Limit: 10, Search: "ORD", Status: "in_transit", ExcludeStatus: "delivered"}) || len(result.Data) != 1 || result.Data[0].ID != "ORD-10246" || result.Pagination.Page != 2 || result.Pagination.TotalPages != 3 {
		t.Fatalf("unexpected list response/filter: %d %+v %+v", response.Code, result, repo.filter)
	}
}
func TestGetByIDNotFound(t *testing.T) {
	app := fiber.New()
	app.Get("/api/shipments/:orderId", NewHandler(application.NewService(&repoFake{}, nil)).GetByID)
	response := request(app, "GET", "/api/shipments/ORD-99999", "")
	var result struct {
		Error string `json:"error"`
	}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if response.Code != 404 || result.Error != "shipment not found" {
		t.Fatalf("unexpected response: %d %+v", response.Code, result)
	}
}
func TestUpdateChangesProvidedFields(t *testing.T) {
	repo := &repoFake{shipment: &domain.Shipment{OrderID: "ORD-10246", Carrier: "Old", Weight: 5, Status: domain.StatusPending}}
	app := fiber.New()
	app.Put("/api/shipments/:orderId", NewHandler(application.NewService(repo, nil)).Update)
	response := request(app, "PUT", "/api/shipments/ORD-10246", `{"carrier":"New"}`)
	if response.Code != 200 || repo.shipment.Carrier != "New" || repo.shipment.Weight != 5 {
		t.Fatalf("partial update lost fields: %d %+v", response.Code, repo.shipment)
	}
}
func TestUpdateStatusCreatesMappedEvent(t *testing.T) {
	repo := &repoFake{shipment: &domain.Shipment{ID: 7, OrderID: "ORD-10246", Destination: "Destination", CurrentCoords: domain.GeoPoint{Lat: 1, Lng: 2}}}
	app := fiber.New()
	app.Patch("/api/shipments/:orderId/status", NewHandler(application.NewService(repo, nil)).UpdateStatus)
	response := request(app, "PATCH", "/api/shipments/ORD-10246/status", `{"status":"in_transit"}`)
	if response.Code != 200 || repo.shipment.Status != domain.StatusInTransit || repo.event == nil || repo.event.Status != "In Transit" || repo.event.Location.Name != "Destination" {
		t.Fatalf("unexpected status update/event: %d %+v %+v", response.Code, repo.shipment, repo.event)
	}
}
func TestDeleteRemovesShipment(t *testing.T) {
	repo := &repoFake{shipment: &domain.Shipment{OrderID: "ORD-10246"}}
	app := fiber.New()
	app.Delete("/api/shipments/:orderId", NewHandler(application.NewService(repo, nil)).Delete)
	response := request(app, "DELETE", "/api/shipments/ORD-10246", "")
	if response.Code != 200 || !repo.deleted {
		t.Fatalf("delete failed: %d deleted=%t", response.Code, repo.deleted)
	}
}

type missingHub struct{}

func (missingHub) FindByID(context.Context, string) (*hubdomain.Hub, error) {
	return nil, domain.ErrNotFound
}

func TestUpdateStatusReturnsInvalidHubIDContract(t *testing.T) {
	repo := &repoFake{shipment: &domain.Shipment{ID: 7, OrderID: "ORD-1"}}
	app := fiber.New()
	app.Patch("/api/shipments/:orderId/status", NewHandler(application.NewService(repo, missingHub{})).UpdateStatus)
	response := request(app, "PATCH", "/api/shipments/ORD-1/status", `{"status":"in_transit","hubId":"HUB-MISSING"}`)
	var result struct {
		Error string `json:"error"`
	}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if response.Code != 400 || result.Error != "invalid hub ID" {
		t.Fatalf("unexpected response: %d %+v", response.Code, result)
	}
}
