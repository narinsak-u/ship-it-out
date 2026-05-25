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
			Name:               "Laem Chabang Port Hub",
			CarrierID:          "THUN",
			Address:            "222 ถ.ท่าเรือแหลมฉบัง ต.แหลมฉบัง อ.ศรีราชา จ.ชลบุรี",
			Coords:             models.GeoPoint{Lat: 13.0833, Lng: 100.8833},
			Capacity:           5000,
			CurrentUtilization: 3400,
			Status:             "active",
			CreatedAt:          time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:                 "HUB-002",
			Name:               "Pattaya Hub",
			CarrierID:          "THUN",
			Address:            "88 ถ.สุขุมวิท ต.หนองปรือ อ.บางละมุง จ.ชลบุรี",
			Coords:             models.GeoPoint{Lat: 12.9236, Lng: 100.8825},
			Capacity:           3000,
			CurrentUtilization: 2100,
			Status:             "active",
			CreatedAt:          time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:                 "HUB-003",
			Name:               "Rayong Hub",
			CarrierID:          "THUN",
			Address:            "150 ถ.สุขุมวิท ต.ท่าประดู่ อ.เมือง จ.ระยอง",
			Coords:             models.GeoPoint{Lat: 12.6814, Lng: 101.2817},
			Capacity:           4000,
			CurrentUtilization: 2800,
			Status:             "active",
			CreatedAt:          time.Date(2026, 2, 20, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:                 "HUB-004",
			Name:               "Chanthaburi Hub",
			CarrierID:          "THUN",
			Address:            "45 ถ.ท่าใหม่ ต.ตลาด อ.เมือง จ.จันทบุรี",
			Coords:             models.GeoPoint{Lat: 12.6096, Lng: 102.1041},
			Capacity:           2000,
			CurrentUtilization: 820,
			Status:             "active",
			CreatedAt:          time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:                 "HUB-005",
			Name:               "Chachoengsao Hub",
			CarrierID:          "THUN",
			Address:            "33 ถ.ฉะเชิงเทรา-บางปะกง ต.หน้าเมือง อ.เมือง จ.ฉะเชิงเทรา",
			Coords:             models.GeoPoint{Lat: 13.6883, Lng: 101.0719},
			Capacity:           2500,
			CurrentUtilization: 210,
			Status:             "active",
			CreatedAt:          time.Date(2026, 4, 5, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:                 "HUB-006",
			Name:               "Trat Hub",
			CarrierID:          "THUN",
			Address:            "12 ถ.ตราด-คลองใหญ่ ต.บางพระ อ.เมือง จ.ตราด",
			Coords:             models.GeoPoint{Lat: 12.2417, Lng: 102.5167},
			Capacity:           1800,
			CurrentUtilization: 990,
			Status:             "maintenance",
			CreatedAt:          time.Date(2026, 4, 15, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, h := range hubs {
		db.Create(&h)
	}
}
