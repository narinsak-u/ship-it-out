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
		HubID             string
		Carrier           string
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
			OrderID: "ORD-10245", TrackingNumber: "TH202600100",
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
			HubID:             "HUB-003",
			Carrier:           "Thun-u-der Express",
			Weight:            "5.2 กก.",
			Items:             2,
			EstimatedDelivery: "May 26, 2026",
			CreatedAt:         "May 22, 2026",
			Progress:          45,
			CurrentCoords:     models.GeoPoint{Lat: 12.85, Lng: 101.5},
			Events: []seedEvent{
				{Timestamp: "May 24, 08:30", Location: models.Location{Name: "Near Ban Bueng", Lat: 12.85, Lng: 101.5}, Status: "In Transit", Description: "Transit to next hub."},
				{Timestamp: "May 22, 14:00", Location: models.Location{Name: "Rayong Hub, 150 ถ.สุขุมวิท ต.ท่าประดู่ อ.เมือง จ.ระยอง", Lat: 12.6814, Lng: 101.2817}, Status: "Departed", Description: "In transit to hub."},
				{Timestamp: "May 22, 09:15", Location: models.Location{Name: "แหลมฉบัง, ศรีราชา, ชลบุรี", Lat: 13.0833, Lng: 100.8833}, Status: "Picked Up", Description: "Parcel collected from sender."},
				{Timestamp: "May 22, 08:00", Location: models.Location{Name: "แหลมฉบัง, ศรีราชา, ชลบุรี", Lat: 13.0833, Lng: 100.8833}, Status: "Label Created", Description: "Awaiting pickup."},
			},
		},
		{
			OrderID: "ORD-10249", TrackingNumber: "TH202600101",
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
			HubID:             "",
			Carrier:           "Thun-u-der Express",
			Weight:            "1.5 กก.",
			Items:             1,
			EstimatedDelivery: "May 27, 2026",
			CreatedAt:         "May 24, 2026",
			Progress:          5,
			CurrentCoords:     models.GeoPoint{Lat: 13.6883, Lng: 101.0719},
			Events: []seedEvent{
				{Timestamp: "May 24, 11:20", Location: models.Location{Name: "หน้าเมือง, เมือง, ฉะเชิงเทรา", Lat: 13.6883, Lng: 101.0719}, Status: "Label Created", Description: "Awaiting pickup."},
			},
		},
		{
			OrderID: "ORD-10250", TrackingNumber: "TH202600102",
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
			Status:            "out_for_delivery",
			HubID:             "HUB-006",
			Carrier:           "Thun-u-der Express",
			Weight:            "3.8 กก.",
			Items:             3,
			EstimatedDelivery: "May 25, 2026",
			CreatedAt:         "May 23, 2026",
			Progress:          75,
			CurrentCoords:     models.GeoPoint{Lat: 12.45, Lng: 101.9},
			Events: []seedEvent{
				{Timestamp: "May 25, 09:00", Location: models.Location{Name: "Trat Hub, 12 ถ.ตราด-คลองใหญ่ ต.บางพระ อ.เมือง จ.ตราด", Lat: 12.2417, Lng: 102.5167}, Status: "Out for Delivery", Description: "Out for delivery."},
				{Timestamp: "May 24, 15:30", Location: models.Location{Name: "Chanthaburi Hub, 45 ถ.ท่าใหม่ ต.ตลาด อ.เมือง จ.จันทบุรี", Lat: 12.6096, Lng: 102.1041}, Status: "In Transit", Description: "Transit to next hub."},
				{Timestamp: "May 23, 16:30", Location: models.Location{Name: "Rayong Hub, 150 ถ.สุขุมวิท ต.ท่าประดู่ อ.เมือง จ.ระยอง", Lat: 12.6814, Lng: 101.2817}, Status: "Departed", Description: "In transit to hub."},
				{Timestamp: "May 23, 13:00", Location: models.Location{Name: "ท่าประดู่, เมือง, ระยอง", Lat: 12.6814, Lng: 101.2817}, Status: "Picked Up", Description: "Parcel collected from sender."},
				{Timestamp: "May 23, 08:00", Location: models.Location{Name: "ท่าประดู่, เมือง, ระยอง", Lat: 12.6814, Lng: 101.2817}, Status: "Label Created", Description: "Awaiting pickup."},
			},
		},
		{
			OrderID: "ORD-10251", TrackingNumber: "TH202600103",
			Customer: models.ContactInfo{
				Name: "สุนิสา ใจดี", Zipcode: "10200",
				SubDistrict: "บางรัก", District: "บางรัก", Province: "กรุงเทพมหานคร",
				Coords: models.GeoPoint{Lat: 13.7279, Lng: 100.5242},
			},
			Receiver: models.ContactInfo{
				Name: "นพดล อินทร์แก้ว", Zipcode: "50000",
				SubDistrict: "ศรีภูมิ", District: "เมือง", Province: "เชียงใหม่",
				Coords: models.GeoPoint{Lat: 18.7883, Lng: 98.9853},
			},
			Status:            "delivered",
			HubID:             "HUB-003",
			Carrier:           "Thun-u-der Express",
			Weight:            "2.1 กก.",
			Items:             1,
			EstimatedDelivery: "May 20, 2026",
			CreatedAt:         "May 16, 2026",
			Progress:          100,
			CurrentCoords:     models.GeoPoint{Lat: 18.7883, Lng: 98.9853},
			Events: []seedEvent{
				{Timestamp: "May 20, 14:30", Location: models.Location{Name: "ศรีภูมิ, เมือง, เชียงใหม่", Lat: 18.7883, Lng: 98.9853}, Status: "Delivered", Description: "Delivered to recipient."},
				{Timestamp: "May 19, 10:00", Location: models.Location{Name: "Rayong Hub, 150 ถ.สุขุมวิท ต.ท่าประดู่ อ.เมือง จ.ระยอง", Lat: 12.6814, Lng: 101.2817}, Status: "Out for Delivery", Description: "Out for delivery."},
				{Timestamp: "May 18, 16:00", Location: models.Location{Name: "Chanthaburi Hub, 45 ถ.ท่าใหม่ ต.ตลาด อ.เมือง จ.จันทบุรี", Lat: 12.6096, Lng: 102.1041}, Status: "In Transit", Description: "Transit to next hub."},
				{Timestamp: "May 17, 09:00", Location: models.Location{Name: "Laem Chabang Port Hub, แหลมฉบัง ศรีราชา ชลบุรี", Lat: 13.0833, Lng: 100.8833}, Status: "Departed", Description: "In transit to hub."},
				{Timestamp: "May 16, 15:00", Location: models.Location{Name: "บางรัก, บางรัก, กรุงเทพมหานคร", Lat: 13.7279, Lng: 100.5242}, Status: "Picked Up", Description: "Parcel collected from sender."},
				{Timestamp: "May 16, 08:00", Location: models.Location{Name: "บางรัก, บางรัก, กรุงเทพมหานคร", Lat: 13.7279, Lng: 100.5242}, Status: "Label Created", Description: "Awaiting pickup."},
			},
		},
		{
			OrderID: "ORD-10252", TrackingNumber: "TH202600104",
			Customer: models.ContactInfo{
				Name: "ธานี ทรัพย์ทวี", Zipcode: "10270",
				SubDistrict: "สำโรงเหนือ", District: "เมือง", Province: "สมุทรปราการ",
				Coords: models.GeoPoint{Lat: 13.6509, Lng: 100.6016},
			},
			Receiver: models.ContactInfo{
				Name: "ไพศาล แซ่เล่า", Zipcode: "83000",
				SubDistrict: "ป่าตอง", District: "กะทู้", Province: "ภูเก็ต",
				Coords: models.GeoPoint{Lat: 7.8961, Lng: 98.2966},
			},
			Status:            "picked_up",
			HubID:             "HUB-001",
			Carrier:           "Thun-u-der Express",
			Weight:            "8.5 กก.",
			Items:             5,
			EstimatedDelivery: "May 29, 2026",
			CreatedAt:         "May 25, 2026",
			Progress:          15,
			CurrentCoords:     models.GeoPoint{Lat: 13.65, Lng: 100.6},
			Events: []seedEvent{
				{Timestamp: "May 25, 14:30", Location: models.Location{Name: "สำโรงเหนือ, เมือง, สมุทรปราการ", Lat: 13.6509, Lng: 100.6016}, Status: "Picked Up", Description: "Parcel collected from sender."},
				{Timestamp: "May 25, 09:00", Location: models.Location{Name: "สำโรงเหนือ, เมือง, สมุทรปราการ", Lat: 13.6509, Lng: 100.6016}, Status: "Label Created", Description: "Awaiting pickup."},
			},
		},
		{
			OrderID: "ORD-10253", TrackingNumber: "TH202600105",
			Customer: models.ContactInfo{
				Name: "รัตนา พิมพ์ทอง", Zipcode: "11000",
				SubDistrict: "ท่าทราย", District: "เมือง", Province: "นนทบุรี",
				Coords: models.GeoPoint{Lat: 13.8548, Lng: 100.5146},
			},
			Receiver: models.ContactInfo{
				Name: "สมศักดิ์ แก้วพล", Zipcode: "57000",
				SubDistrict: "เวียง", District: "เมือง", Province: "เชียงราย",
				Coords: models.GeoPoint{Lat: 19.9072, Lng: 99.8325},
			},
			Status:            "in_transit",
			HubID:             "HUB-002",
			Carrier:           "Thun-u-der Express",
			Weight:            "4.2 กก.",
			Items:             2,
			EstimatedDelivery: "May 28, 2026",
			CreatedAt:         "May 24, 2026",
			Progress:          35,
			CurrentCoords:     models.GeoPoint{Lat: 16.5, Lng: 100.5},
			Events: []seedEvent{
				{Timestamp: "May 25, 11:00", Location: models.Location{Name: "Nakhon Sawan", Lat: 15.7167, Lng: 100.1333}, Status: "In Transit", Description: "Transit to next hub."},
				{Timestamp: "May 24, 16:00", Location: models.Location{Name: "Pattaya Hub, พัทยา บางละมุง ชลบุรี", Lat: 12.9236, Lng: 100.8825}, Status: "Departed", Description: "In transit to hub."},
				{Timestamp: "May 24, 10:00", Location: models.Location{Name: "ท่าทราย, เมือง, นนทบุรี", Lat: 13.8548, Lng: 100.5146}, Status: "Picked Up", Description: "Parcel collected from sender."},
				{Timestamp: "May 24, 08:00", Location: models.Location{Name: "ท่าทราย, เมือง, นนทบุรี", Lat: 13.8548, Lng: 100.5146}, Status: "Label Created", Description: "Awaiting pickup."},
			},
		},
		{
			OrderID: "ORD-10254", TrackingNumber: "TH202600106",
			Customer: models.ContactInfo{
				Name: "สุธีร์ ศรีวัฒน์", Zipcode: "12000",
				SubDistrict: "ประชาธิปัตย์", District: "ธัญบุรี", Province: "ปทุมธานี",
				Coords: models.GeoPoint{Lat: 13.9782, Lng: 100.6147},
			},
			Receiver: models.ContactInfo{
				Name: "อับดุลเลาะ สะมะแอ", Zipcode: "90110",
				SubDistrict: "คอหงส์", District: "หาดใหญ่", Province: "สงขลา",
				Coords: models.GeoPoint{Lat: 7.0088, Lng: 100.4747},
			},
			Status:            "delayed",
			HubID:             "HUB-004",
			Carrier:           "Thun-u-der Express",
			Weight:            "6.7 กก.",
			Items:             4,
			EstimatedDelivery: "May 24, 2026",
			CreatedAt:         "May 20, 2026",
			Progress:          60,
			CurrentCoords:     models.GeoPoint{Lat: 11.5, Lng: 99.5},
			Events: []seedEvent{
				{Timestamp: "May 25, 08:00", Location: models.Location{Name: "Chanthaburi Hub, 45 ถ.ท่าใหม่ ต.ตลาด อ.เมือง จ.จันทบุรี", Lat: 12.6096, Lng: 102.1041}, Status: "Delayed", Description: "Unexpected issue encountered."},
				{Timestamp: "May 23, 14:00", Location: models.Location{Name: "Chanthaburi Hub, 45 ถ.ท่าใหม่ ต.ตลาด อ.เมือง จ.จันทบุรี", Lat: 12.6096, Lng: 102.1041}, Status: "In Transit", Description: "Transit to next hub."},
				{Timestamp: "May 22, 09:00", Location: models.Location{Name: "Laem Chabang Port Hub, แหลมฉบัง ศรีราชา ชลบุรี", Lat: 13.0833, Lng: 100.8833}, Status: "Departed", Description: "In transit to hub."},
				{Timestamp: "May 21, 16:00", Location: models.Location{Name: "ประชาธิปัตย์, ธัญบุรี, ปทุมธานี", Lat: 13.9782, Lng: 100.6147}, Status: "Picked Up", Description: "Parcel collected from sender."},
				{Timestamp: "May 21, 08:00", Location: models.Location{Name: "ประชาธิปัตย์, ธัญบุรี, ปทุมธานี", Lat: 13.9782, Lng: 100.6147}, Status: "Label Created", Description: "Awaiting pickup."},
			},
		},
		{
			OrderID: "ORD-10255", TrackingNumber: "TH202600107",
			Customer: models.ContactInfo{
				Name: "กาญจนา ศิริโชค", Zipcode: "13000",
				SubDistrict: "ประตูชัย", District: "พระนครศรีอยุธยา", Province: "พระนครศรีอยุธยา",
				Coords: models.GeoPoint{Lat: 14.3542, Lng: 100.5547},
			},
			Receiver: models.ContactInfo{
				Name: "ธนพล จันทร์ศรี", Zipcode: "40000",
				SubDistrict: "ในเมือง", District: "เมือง", Province: "ขอนแก่น",
				Coords: models.GeoPoint{Lat: 16.4322, Lng: 102.8236},
			},
			Status:            "pending",
			HubID:             "",
			Carrier:           "Thun-u-der Express",
			Weight:            "3.3 กก.",
			Items:             2,
			EstimatedDelivery: "May 29, 2026",
			CreatedAt:         "May 26, 2026",
			Progress:          0,
			CurrentCoords:     models.GeoPoint{Lat: 14.3542, Lng: 100.5547},
			Events: []seedEvent{
				{Timestamp: "May 26, 10:00", Location: models.Location{Name: "ประตูชัย, พระนครศรีอยุธยา, พระนครศรีอยุธยา", Lat: 14.3542, Lng: 100.5547}, Status: "Label Created", Description: "Awaiting pickup."},
			},
		},
		{
			OrderID: "ORD-10256", TrackingNumber: "TH202600108",
			Customer: models.ContactInfo{
				Name: "นฤมล สุขเกษม", Zipcode: "10100",
				SubDistrict: "ปทุมวัน", District: "ปทุมวัน", Province: "กรุงเทพมหานคร",
				Coords: models.GeoPoint{Lat: 13.7466, Lng: 100.5326},
			},
			Receiver: models.ContactInfo{
				Name: "อดุลย์ กล้าหาญ", Zipcode: "30000",
				SubDistrict: "ในเมือง", District: "เมือง", Province: "นครราชสีมา",
				Coords: models.GeoPoint{Lat: 14.975, Lng: 102.0825},
			},
			Status:            "delivered",
			HubID:             "HUB-001",
			Carrier:           "Thun-u-der Express",
			Weight:            "1.8 กก.",
			Items:             1,
			EstimatedDelivery: "May 23, 2026",
			CreatedAt:         "May 19, 2026",
			Progress:          100,
			CurrentCoords:     models.GeoPoint{Lat: 14.975, Lng: 102.0825},
			Events: []seedEvent{
				{Timestamp: "May 23, 11:30", Location: models.Location{Name: "ในเมือง, เมือง, นครราชสีมา", Lat: 14.975, Lng: 102.0825}, Status: "Delivered", Description: "Delivered to recipient."},
				{Timestamp: "May 22, 15:00", Location: models.Location{Name: "Laem Chabang Port Hub, แหลมฉบัง ศรีราชา ชลบุรี", Lat: 13.0833, Lng: 100.8833}, Status: "Out for Delivery", Description: "Out for delivery."},
				{Timestamp: "May 21, 10:00", Location: models.Location{Name: "Chachoengsao Hub, หน้าเมือง เมือง ฉะเชิงเทรา", Lat: 13.6883, Lng: 101.0719}, Status: "In Transit", Description: "Transit to next hub."},
				{Timestamp: "May 20, 14:00", Location: models.Location{Name: "Laem Chabang Port Hub, แหลมฉบัง ศรีราชา ชลบุรี", Lat: 13.0833, Lng: 100.8833}, Status: "Departed", Description: "In transit to hub."},
				{Timestamp: "May 19, 16:30", Location: models.Location{Name: "ปทุมวัน, ปทุมวัน, กรุงเทพมหานคร", Lat: 13.7466, Lng: 100.5326}, Status: "Picked Up", Description: "Parcel collected from sender."},
				{Timestamp: "May 19, 08:00", Location: models.Location{Name: "ปทุมวัน, ปทุมวัน, กรุงเทพมหานคร", Lat: 13.7466, Lng: 100.5326}, Status: "Label Created", Description: "Awaiting pickup."},
			},
		},
		{
			OrderID: "ORD-10257", TrackingNumber: "TH202600109",
			Customer: models.ContactInfo{
				Name: "สาโรจน์ เจริญสุข", Zipcode: "20110",
				SubDistrict: "บางพระ", District: "ศรีราชา", Province: "ชลบุรี",
				Coords: models.GeoPoint{Lat: 13.1173, Lng: 100.9256},
			},
			Receiver: models.ContactInfo{
				Name: "วาสนา คงมั่น", Zipcode: "84000",
				SubDistrict: "ตลาด", District: "เมือง", Province: "สุราษฎร์ธานี",
				Coords: models.GeoPoint{Lat: 9.1382, Lng: 99.3214},
			},
			Status:            "out_for_delivery",
			HubID:             "HUB-005",
			Carrier:           "Thun-u-der Express",
			Weight:            "10.0 กก.",
			Items:             6,
			EstimatedDelivery: "May 26, 2026",
			CreatedAt:         "May 22, 2026",
			Progress:          80,
			CurrentCoords:     models.GeoPoint{Lat: 9.5, Lng: 99.2},
			Events: []seedEvent{
				{Timestamp: "May 26, 09:00", Location: models.Location{Name: "Chachoengsao Hub, หน้าเมือง เมือง ฉะเชิงเทรา", Lat: 13.6883, Lng: 101.0719}, Status: "Out for Delivery", Description: "Out for delivery."},
				{Timestamp: "May 25, 14:00", Location: models.Location{Name: "Chachoengsao Hub, หน้าเมือง เมือง ฉะเชิงเทรา", Lat: 13.6883, Lng: 101.0719}, Status: "In Transit", Description: "Transit to next hub."},
				{Timestamp: "May 24, 10:00", Location: models.Location{Name: "Laem Chabang Port Hub, แหลมฉบัง ศรีราชา ชลบุรี", Lat: 13.0833, Lng: 100.8833}, Status: "Departed", Description: "In transit to hub."},
				{Timestamp: "May 23, 14:00", Location: models.Location{Name: "บางพระ, ศรีราชา, ชลบุรี", Lat: 13.1173, Lng: 100.9256}, Status: "Picked Up", Description: "Parcel collected from sender."},
				{Timestamp: "May 23, 08:00", Location: models.Location{Name: "บางพระ, ศรีราชา, ชลบุรี", Lat: 13.1173, Lng: 100.9256}, Status: "Label Created", Description: "Awaiting pickup."},
			},
		},
		{
			OrderID: "ORD-10258", TrackingNumber: "TH202600110",
			Customer: models.ContactInfo{
				Name: "ประเสริฐ วงศ์ดี", Zipcode: "21160",
				SubDistrict: "นิคมพัฒนา", District: "นิคมพัฒนา", Province: "ระยอง",
				Coords: models.GeoPoint{Lat: 12.7167, Lng: 101.15},
			},
			Receiver: models.ContactInfo{
				Name: "บุญธรรม พิมพา", Zipcode: "41000",
				SubDistrict: "หมากแข้ง", District: "เมือง", Province: "อุดรธานี",
				Coords: models.GeoPoint{Lat: 17.4132, Lng: 102.7856},
			},
			Status:            "in_transit",
			HubID:             "HUB-003",
			Carrier:           "Thun-u-der Express",
			Weight:            "7.1 กก.",
			Items:             3,
			EstimatedDelivery: "May 28, 2026",
			CreatedAt:         "May 24, 2026",
			Progress:          40,
			CurrentCoords:     models.GeoPoint{Lat: 15.0, Lng: 102.0},
			Events: []seedEvent{
				{Timestamp: "May 25, 16:00", Location: models.Location{Name: "Nakhon Ratchasima", Lat: 14.975, Lng: 102.0825}, Status: "In Transit", Description: "Transit to next hub."},
				{Timestamp: "May 24, 18:00", Location: models.Location{Name: "Rayong Hub, 150 ถ.สุขุมวิท ต.ท่าประดู่ อ.เมือง จ.ระยอง", Lat: 12.6814, Lng: 101.2817}, Status: "Departed", Description: "In transit to hub."},
				{Timestamp: "May 24, 11:00", Location: models.Location{Name: "นิคมพัฒนา, นิคมพัฒนา, ระยอง", Lat: 12.7167, Lng: 101.15}, Status: "Picked Up", Description: "Parcel collected from sender."},
				{Timestamp: "May 24, 08:00", Location: models.Location{Name: "นิคมพัฒนา, นิคมพัฒนา, ระยอง", Lat: 12.7167, Lng: 101.15}, Status: "Label Created", Description: "Awaiting pickup."},
			},
		},
		{
			OrderID: "ORD-10259", TrackingNumber: "TH202600111",
			Customer: models.ContactInfo{
				Name: "สุดา มากสุข", Zipcode: "24110",
				SubDistrict: "บางคล้า", District: "บางคล้า", Province: "ฉะเชิงเทรา",
				Coords: models.GeoPoint{Lat: 13.7167, Lng: 101.2},
			},
			Receiver: models.ContactInfo{
				Name: "สมพร รักษ์ดี", Zipcode: "80000",
				SubDistrict: "ท่าวัง", District: "เมือง", Province: "นครศรีธรรมราช",
				Coords: models.GeoPoint{Lat: 8.4333, Lng: 99.9667},
			},
			Status:            "departed",
			HubID:             "HUB-005",
			Carrier:           "Thun-u-der Express",
			Weight:            "5.5 กก.",
			Items:             3,
			EstimatedDelivery: "May 28, 2026",
			CreatedAt:         "May 25, 2026",
			Progress:          20,
			CurrentCoords:     models.GeoPoint{Lat: 13.7, Lng: 101.1},
			Events: []seedEvent{
				{Timestamp: "May 25, 17:00", Location: models.Location{Name: "Chachoengsao Hub, หน้าเมือง เมือง ฉะเชิงเทรา", Lat: 13.6883, Lng: 101.0719}, Status: "Departed", Description: "In transit to hub."},
				{Timestamp: "May 25, 13:00", Location: models.Location{Name: "บางคล้า, บางคล้า, ฉะเชิงเทรา", Lat: 13.7167, Lng: 101.2}, Status: "Picked Up", Description: "Parcel collected from sender."},
				{Timestamp: "May 25, 08:00", Location: models.Location{Name: "บางคล้า, บางคล้า, ฉะเชิงเทรา", Lat: 13.7167, Lng: 101.2}, Status: "Label Created", Description: "Awaiting pickup."},
			},
		},
		{
			OrderID: "ORD-10260", TrackingNumber: "TH202600112",
			Customer: models.ContactInfo{
				Name: "พิมล สุขใจ", Zipcode: "10250",
				SubDistrict: "คลองเตย", District: "คลองเตย", Province: "กรุงเทพมหานคร",
				Coords: models.GeoPoint{Lat: 13.7122, Lng: 100.5638},
			},
			Receiver: models.ContactInfo{
				Name: "วีระ วงค์คำ", Zipcode: "34000",
				SubDistrict: "ในเมือง", District: "เมือง", Province: "อุบลราชธานี",
				Coords: models.GeoPoint{Lat: 15.2296, Lng: 104.8603},
			},
			Status:            "delivered",
			HubID:             "HUB-004",
			Carrier:           "Thun-u-der Express",
			Weight:            "4.5 กก.",
			Items:             2,
			EstimatedDelivery: "May 24, 2026",
			CreatedAt:         "May 20, 2026",
			Progress:          100,
			CurrentCoords:     models.GeoPoint{Lat: 15.2296, Lng: 104.8603},
			Events: []seedEvent{
				{Timestamp: "May 24, 10:00", Location: models.Location{Name: "ในเมือง, เมือง, อุบลราชธานี", Lat: 15.2296, Lng: 104.8603}, Status: "Delivered", Description: "Delivered to recipient."},
				{Timestamp: "May 23, 14:00", Location: models.Location{Name: "Chanthaburi Hub, 45 ถ.ท่าใหม่ ต.ตลาด อ.เมือง จ.จันทบุรี", Lat: 12.6096, Lng: 102.1041}, Status: "Out for Delivery", Description: "Out for delivery."},
				{Timestamp: "May 22, 16:00", Location: models.Location{Name: "Chanthaburi Hub, 45 ถ.ท่าใหม่ ต.ตลาด อ.เมือง จ.จันทบุรี", Lat: 12.6096, Lng: 102.1041}, Status: "In Transit", Description: "Transit to next hub."},
				{Timestamp: "May 21, 14:00", Location: models.Location{Name: "Laem Chabang Port Hub, แหลมฉบัง ศรีราชา ชลบุรี", Lat: 13.0833, Lng: 100.8833}, Status: "Departed", Description: "In transit to hub."},
				{Timestamp: "May 20, 17:00", Location: models.Location{Name: "คลองเตย, คลองเตย, กรุงเทพมหานคร", Lat: 13.7122, Lng: 100.5638}, Status: "Picked Up", Description: "Parcel collected from sender."},
				{Timestamp: "May 20, 08:00", Location: models.Location{Name: "คลองเตย, คลองเตย, กรุงเทพมหานคร", Lat: 13.7122, Lng: 100.5638}, Status: "Label Created", Description: "Awaiting pickup."},
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
			HubID:             s.HubID,
			Carrier:           s.Carrier,
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
