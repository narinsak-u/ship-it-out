package seed

import (
	"time"

	"github.com/narinsak-u/backend/internal/models"
	"gorm.io/gorm"
)

func SeedShipments(db *gorm.DB) {
	var count int64
	db.Model(&models.Shipment{}).Count(&count)
	if count > 0 {
		return
	}

	parse := func(layout, value string) time.Time {
		t, err := time.Parse(layout, value)
		if err != nil {
			panic("seed: failed to parse time: " + err.Error())
		}
		return t
	}

	type seedEvent struct {
		Timestamp   string
		Location    models.Location
		Status      string
		Description string
	}

	type seedShipment struct {
		OrderID           string
		TrackingNumber    string
		Customer          models.ContactInfo
		Receiver          models.ContactInfo
		Status            string
		Weight            string
		Items             int
		EstimatedDelivery string
		CreatedAt         string
		Progress          float64
		CurrentCoords     models.GeoPoint
		Events            []seedEvent
	}

	seeds := []seedShipment{
		{
			OrderID: "ORD-10245", TrackingNumber: "TRK-9F2A-44B1",
			Customer: models.ContactInfo{
				Name: "สมชาย วงศ์เจริญ", Zipcode: "20110",
				SubDistrict: "แหลมฉบัง", District: "ศรีราชา", Province: "ชลบุรี",
				Coords: models.GeoPoint{Lat: 13.0833, Lng: 100.8833},
			},
			Receiver: models.ContactInfo{
				Name: "มาลี ทองดี", Zipcode: "22000",
				SubDistrict: "จันทนิมิต", District: "เมือง", Province: "จันทบุรี",
				Coords: models.GeoPoint{Lat: 12.6096, Lng: 102.1041},
			},
			Status:            "in_transit",
			Weight:            "5.2 กก.",
			Items:             2,
			EstimatedDelivery: "May 26, 2026",
			CreatedAt:         "May 22, 2026",
			Progress:          45,
			CurrentCoords:     models.GeoPoint{Lat: 12.85, Lng: 101.5},
			Events: []seedEvent{
				{Timestamp: "May 24, 08:30", Location: models.Location{Name: "Near Ban Bueng", Lat: 12.85, Lng: 101.5}, Status: "In transit", Description: "กำลังเดินทางผ่าน อ.บ้านบึง"},
				{Timestamp: "May 22, 14:00", Location: models.Location{Name: "Laem Chabang Port Hub", Lat: 13.0833, Lng: 100.8833}, Status: "Departed", Description: "ออกจากคลังสินค้าแหลมฉบัง"},
				{Timestamp: "May 22, 09:15", Location: models.Location{Name: "Laem Chabang Port Hub", Lat: 13.0833, Lng: 100.8833}, Status: "Picked up", Description: "รับพัสดุจากผู้ส่ง"},
			},
		},
		{
			OrderID: "ORD-10249", TrackingNumber: "TRK-FF02-1188",
			Customer: models.ContactInfo{
				Name: "วิมล ศรีสุวรรณ", Zipcode: "24000",
				SubDistrict: "หน้าเมือง", District: "เมือง", Province: "ฉะเชิงเทรา",
				Coords: models.GeoPoint{Lat: 13.6883, Lng: 101.0719},
			},
			Receiver: models.ContactInfo{
				Name: "กิตติพงศ์ แก้ววิเศษ", Zipcode: "20150",
				SubDistrict: "หนองปรือ", District: "บางละมุง", Province: "ชลบุรี",
				Coords: models.GeoPoint{Lat: 12.9236, Lng: 100.8825},
			},
			Status:            "pending",
			Weight:            "1.5 กก.",
			Items:             1,
			EstimatedDelivery: "May 27, 2026",
			CreatedAt:         "May 24, 2026",
			Progress:          5,
			CurrentCoords:     models.GeoPoint{Lat: 13.6883, Lng: 101.0719},
			Events: []seedEvent{
				{Timestamp: "May 24, 11:20", Location: models.Location{Name: "Chachoengsao Hub", Lat: 13.6883, Lng: 101.0719}, Status: "Label created", Description: "รอรับพัสดุจากผู้ส่ง"},
			},
		},
		{
			OrderID: "ORD-10250", TrackingNumber: "TRK-5E73-220B",
			Customer: models.ContactInfo{
				Name: "วิชัย สมบูรณ์", Zipcode: "21000",
				SubDistrict: "ท่าประดู่", District: "เมือง", Province: "ระยอง",
				Coords: models.GeoPoint{Lat: 12.6814, Lng: 101.2817},
			},
			Receiver: models.ContactInfo{
				Name: "ประภาสิริ วัฒนา", Zipcode: "23000",
				SubDistrict: "บางพระ", District: "เมือง", Province: "ตราด",
				Coords: models.GeoPoint{Lat: 12.2417, Lng: 102.5167},
			},
			Status:            "in_transit",
			Weight:            "3.8 กก.",
			Items:             3,
			EstimatedDelivery: "May 25, 2026",
			CreatedAt:         "May 23, 2026",
			Progress:          60,
			CurrentCoords:     models.GeoPoint{Lat: 12.45, Lng: 101.9},
			Events: []seedEvent{
				{Timestamp: "May 24, 06:45", Location: models.Location{Name: "Near Klaeng", Lat: 12.45, Lng: 101.9}, Status: "In transit", Description: "กำลังเดินทางผ่าน อ.แกลง"},
				{Timestamp: "May 23, 16:30", Location: models.Location{Name: "Rayong Hub", Lat: 12.6814, Lng: 101.2817}, Status: "Departed", Description: "ออกจากคลังสินค้าระยอง"},
				{Timestamp: "May 23, 13:00", Location: models.Location{Name: "Rayong Hub", Lat: 12.6814, Lng: 101.2817}, Status: "Picked up", Description: "รับพัสดุจากผู้ส่ง"},
			},
		},
	}

	composeAddress := func(c models.ContactInfo) string {
		return c.SubDistrict + ", " + c.District + ", " + c.Province
	}

	for _, s := range seeds {
		shipment := models.Shipment{
			OrderID:           s.OrderID,
			TrackingNumber:    s.TrackingNumber,
			Customer:          s.Customer,
			Receiver:          s.Receiver,
			Origin:            composeAddress(s.Customer),
			Destination:       composeAddress(s.Receiver),
			Status:            s.Status,
			Carrier:           "Thun-u-der Express",
			Weight:            s.Weight,
			Items:             s.Items,
			EstimatedDelivery: parse("January 2, 2006", s.EstimatedDelivery),
			CreatedAt:         parse("January 2, 2006", s.CreatedAt),
			Progress:          s.Progress,
			CurrentCoords:     s.CurrentCoords,
		}
		db.Create(&shipment)

		for _, e := range s.Events {
			event := models.ShipmentEvent{
				ShipmentID:  shipment.ID,
				Status:      e.Status,
				Location:    e.Location,
				Description: e.Description,
				CreatedAt:   parse("January 2, 15:04", e.Timestamp),
			}
			db.Create(&event)
		}
	}
}
