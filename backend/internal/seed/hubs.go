package seed

import (
	"time"

	"github.com/narinsak-u/backend/internal/models"
	"gorm.io/gorm"
)

func SeedHubs(db *gorm.DB) {
	var count int64
	db.Model(&models.Hub{}).Count(&count)
	if count > 0 {
		return
	}

	hubs := []models.Hub{
		{
			ID:                 "HUB-001",
			Name:               "Bangkok Hub",
			CarrierID:          "THUN",
			Address:            "Bangkok",
			Coords:             models.GeoPoint{Lat: 13.7563, Lng: 100.5018},
			Capacity:           5000,
			CurrentUtilization: 3400,
			Status:             "active",
			CreatedAt:          time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:                 "HUB-002",
			Name:               "Chonburi Hub",
			CarrierID:          "THUN",
			Address:            "Chonburi",
			Coords:             models.GeoPoint{Lat: 13.3611, Lng: 100.9847},
			Capacity:           3000,
			CurrentUtilization: 2100,
			Status:             "active",
			CreatedAt:          time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:                 "HUB-003",
			Name:               "Kanchanaburi Hub",
			CarrierID:          "THUN",
			Address:            "Kanchanaburi",
			Coords:             models.GeoPoint{Lat: 14.0199, Lng: 99.4778},
			Capacity:           4000,
			CurrentUtilization: 2800,
			Status:             "active",
			CreatedAt:          time.Date(2026, 2, 20, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:                 "HUB-004",
			Name:               "Chiang Mai Hub",
			CarrierID:          "THUN",
			Address:            "Chiang Mai",
			Coords:             models.GeoPoint{Lat: 18.7883, Lng: 98.9853},
			Capacity:           2000,
			CurrentUtilization: 820,
			Status:             "active",
			CreatedAt:          time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:                 "HUB-005",
			Name:               "Phuket Hub",
			CarrierID:          "THUN",
			Address:            "Phuket",
			Coords:             models.GeoPoint{Lat: 7.8804, Lng: 98.3923},
			Capacity:           2500,
			CurrentUtilization: 210,
			Status:             "active",
			CreatedAt:          time.Date(2026, 4, 5, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:                 "HUB-006",
			Name:               "Korat Hub",
			CarrierID:          "THUN",
			Address:            "Nakhon Ratchasima",
			Coords:             models.GeoPoint{Lat: 14.9799, Lng: 102.0977},
			Capacity:           3500,
			CurrentUtilization: 1800,
			Status:             "active",
			CreatedAt:          time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:                 "HUB-007",
			Name:               "Khon Kaen Hub",
			CarrierID:          "THUN",
			Address:            "Khon Kaen",
			Coords:             models.GeoPoint{Lat: 16.4322, Lng: 102.8236},
			Capacity:           2800,
			CurrentUtilization: 1400,
			Status:             "active",
			CreatedAt:          time.Date(2026, 5, 10, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:                 "HUB-008",
			Name:               "Udon Thani Hub",
			CarrierID:          "THUN",
			Address:            "Udon Thani",
			Coords:             models.GeoPoint{Lat: 17.3647, Lng: 102.8159},
			Capacity:           2200,
			CurrentUtilization: 1100,
			Status:             "active",
			CreatedAt:          time.Date(2026, 5, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:                 "HUB-009",
			Name:               "Ubon Ratchathani Hub",
			CarrierID:          "THUN",
			Address:            "Ubon Ratchathani",
			Coords:             models.GeoPoint{Lat: 15.2448, Lng: 104.8474},
			Capacity:           2000,
			CurrentUtilization: 900,
			Status:             "active",
			CreatedAt:          time.Date(2026, 5, 20, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:                 "HUB-010",
			Name:               "Buriram Hub",
			CarrierID:          "THUN",
			Address:            "Buriram",
			Coords:             models.GeoPoint{Lat: 14.9951, Lng: 103.1035},
			Capacity:           1800,
			CurrentUtilization: 600,
			Status:             "active",
			CreatedAt:          time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, h := range hubs {
		db.Create(&h)
	}
}
