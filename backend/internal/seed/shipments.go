package seed

import (
	"math/rand"
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

	rng := rand.New(rand.NewSource(42))

	type person struct {
		Name        string
		Zipcode     string
		SubDistrict string
		District    string
		Province    string
		Lat         float64
		Lng         float64
	}

	customers := []person{
		{Name: "สมชาย วงศ์เจริญ", Zipcode: "20110", SubDistrict: "แหลมฉบัง", District: "ศรีราชา", Province: "ชลบุรี", Lat: 13.0833, Lng: 100.8833},
		{Name: "วิมล ศรีสุวรรณ", Zipcode: "24000", SubDistrict: "หน้าเมือง", District: "เมือง", Province: "ฉะเชิงเทรา", Lat: 13.6883, Lng: 101.0719},
		{Name: "วิชัย สมบูรณ์", Zipcode: "21000", SubDistrict: "ท่าประดู่", District: "เมือง", Province: "ระยอง", Lat: 12.6814, Lng: 101.2817},
		{Name: "สุนิสา ใจดี", Zipcode: "10200", SubDistrict: "บางรัก", District: "บางรัก", Province: "กรุงเทพมหานคร", Lat: 13.7279, Lng: 100.5242},
		{Name: "ธานี ทรัพย์ทวี", Zipcode: "10270", SubDistrict: "สำโรงเหนือ", District: "เมือง", Province: "สมุทรปราการ", Lat: 13.6509, Lng: 100.6016},
		{Name: "รัตนา พิมพ์ทอง", Zipcode: "11000", SubDistrict: "ท่าทราย", District: "เมือง", Province: "นนทบุรี", Lat: 13.8548, Lng: 100.5146},
		{Name: "สุธีร์ ศรีวัฒน์", Zipcode: "12000", SubDistrict: "ประชาธิปัตย์", District: "ธัญบุรี", Province: "ปทุมธานี", Lat: 13.9782, Lng: 100.6147},
		{Name: "กาญจนา ศิริโชค", Zipcode: "13000", SubDistrict: "ประตูชัย", District: "พระนครศรีอยุธยา", Province: "พระนครศรีอยุธยา", Lat: 14.3542, Lng: 100.5547},
		{Name: "นฤมล สุขเกษม", Zipcode: "10100", SubDistrict: "ปทุมวัน", District: "ปทุมวัน", Province: "กรุงเทพมหานคร", Lat: 13.7466, Lng: 100.5326},
		{Name: "สาโรจน์ เจริญสุข", Zipcode: "20110", SubDistrict: "บางพระ", District: "ศรีราชา", Province: "ชลบุรี", Lat: 13.1173, Lng: 100.9256},
		{Name: "ประเสริฐ วงศ์ดี", Zipcode: "21160", SubDistrict: "นิคมพัฒนา", District: "นิคมพัฒนา", Province: "ระยอง", Lat: 12.7167, Lng: 101.15},
		{Name: "สุดา มากสุข", Zipcode: "24110", SubDistrict: "บางคล้า", District: "บางคล้า", Province: "ฉะเชิงเทรา", Lat: 13.7167, Lng: 101.2},
		{Name: "พิมล สุขใจ", Zipcode: "10250", SubDistrict: "คลองเตย", District: "คลองเตย", Province: "กรุงเทพมหานคร", Lat: 13.7122, Lng: 100.5638},
		{Name: "อรัญญา มั่นคง", Zipcode: "22000", SubDistrict: "จันทนิมิต", District: "เมือง", Province: "จันทบุรี", Lat: 12.6096, Lng: 102.1041},
		{Name: "ปรีชา ทรัพย์เจริญ", Zipcode: "20150", SubDistrict: "หนองปรือ", District: "บางละมุง", Province: "ชลบุรี", Lat: 12.9236, Lng: 100.8825},
		{Name: "ดารา รักไทย", Zipcode: "23000", SubDistrict: "บางพระ", District: "เมือง", Province: "ตราด", Lat: 12.2417, Lng: 102.5167},
		{Name: "มานพ คงสมบูรณ์", Zipcode: "84000", SubDistrict: "ตลาด", District: "เมือง", Province: "สุราษฎร์ธานี", Lat: 9.1382, Lng: 99.3214},
		{Name: "จินตนา แก้วประเสริฐ", Zipcode: "21000", SubDistrict: "ท่าประดู่", District: "เมือง", Province: "ระยอง", Lat: 12.6814, Lng: 101.2817},
		{Name: "สมบูรณ์ เจริญผล", Zipcode: "13000", SubDistrict: "ประตูชัย", District: "พระนครศรีอยุธยา", Province: "พระนครศรีอยุธยา", Lat: 14.3542, Lng: 100.5547},
		{Name: "วรรณา ศิริสวัสดิ์", Zipcode: "24000", SubDistrict: "หน้าเมือง", District: "เมือง", Province: "ฉะเชิงเทรา", Lat: 13.6883, Lng: 101.0719},
	}

	receivers := []person{
		{Name: "มาลี ทองดี", Zipcode: "22000", SubDistrict: "จันทนิมิต", District: "เมือง", Province: "จันทบุรี", Lat: 12.6096, Lng: 102.1041},
		{Name: "กิตติพงศ์ แก้ววิเศษ", Zipcode: "20150", SubDistrict: "หนองปรือ", District: "บางละมุง", Province: "ชลบุรี", Lat: 12.9236, Lng: 100.8825},
		{Name: "ประภาสิริ วัฒนา", Zipcode: "23000", SubDistrict: "บางพระ", District: "เมือง", Province: "ตราด", Lat: 12.2417, Lng: 102.5167},
		{Name: "นพดล อินทร์แก้ว", Zipcode: "50000", SubDistrict: "ศรีภูมิ", District: "เมือง", Province: "เชียงใหม่", Lat: 18.7883, Lng: 98.9853},
		{Name: "ไพศาล แซ่เล่า", Zipcode: "83000", SubDistrict: "ป่าตอง", District: "กะทู้", Province: "ภูเก็ต", Lat: 7.8961, Lng: 98.2966},
		{Name: "สมศักดิ์ แก้วพล", Zipcode: "57000", SubDistrict: "เวียง", District: "เมือง", Province: "เชียงราย", Lat: 19.9072, Lng: 99.8325},
		{Name: "อับดุลเลาะ สะมะแอ", Zipcode: "90110", SubDistrict: "คอหงส์", District: "หาดใหญ่", Province: "สงขลา", Lat: 7.0088, Lng: 100.4747},
		{Name: "ธนพล จันทร์ศรี", Zipcode: "40000", SubDistrict: "ในเมือง", District: "เมือง", Province: "ขอนแก่น", Lat: 16.4322, Lng: 102.8236},
		{Name: "อดุลย์ กล้าหาญ", Zipcode: "30000", SubDistrict: "ในเมือง", District: "เมือง", Province: "นครราชสีมา", Lat: 14.975, Lng: 102.0825},
		{Name: "วาสนา คงมั่น", Zipcode: "84000", SubDistrict: "ตลาด", District: "เมือง", Province: "สุราษฎร์ธานี", Lat: 9.1382, Lng: 99.3214},
		{Name: "บุญธรรม พิมพา", Zipcode: "41000", SubDistrict: "หมากแข้ง", District: "เมือง", Province: "อุดรธานี", Lat: 17.4132, Lng: 102.7856},
		{Name: "สมพร รักษ์ดี", Zipcode: "80000", SubDistrict: "ท่าวัง", District: "เมือง", Province: "นครศรีธรรมราช", Lat: 8.4333, Lng: 99.9667},
		{Name: "วีระ วงค์คำ", Zipcode: "34000", SubDistrict: "ในเมือง", District: "เมือง", Province: "อุบลราชธานี", Lat: 15.2296, Lng: 104.8603},
		{Name: "ธัญญา เรืองศรี", Zipcode: "22000", SubDistrict: "จันทนิมิต", District: "เมือง", Province: "จันทบุรี", Lat: 12.6096, Lng: 102.1041},
		{Name: "สุชาติ พัฒนา", Zipcode: "20110", SubDistrict: "บางพระ", District: "ศรีราชา", Province: "ชลบุรี", Lat: 13.1173, Lng: 100.9256},
		{Name: "บุญมี ชัยวัฒน์", Zipcode: "50000", SubDistrict: "ศรีภูมิ", District: "เมือง", Province: "เชียงใหม่", Lat: 18.7883, Lng: 98.9853},
		{Name: "สมหมาย ใจบุญ", Zipcode: "40000", SubDistrict: "ในเมือง", District: "เมือง", Province: "ขอนแก่น", Lat: 16.4322, Lng: 102.8236},
		{Name: "ประทุม ทองสุข", Zipcode: "30000", SubDistrict: "ในเมือง", District: "เมือง", Province: "นครราชสีมา", Lat: 14.975, Lng: 102.0825},
		{Name: "ส้มจีน แซ่ตั้ง", Zipcode: "90110", SubDistrict: "คอหงส์", District: "หาดใหญ่", Province: "สงขลา", Lat: 7.0088, Lng: 100.4747},
		{Name: "องอาจ เดชะ", Zipcode: "83000", SubDistrict: "ป่าตอง", District: "กะทู้", Province: "ภูเก็ต", Lat: 7.8961, Lng: 98.2966},
	}

	hubs := []string{"HUB-001", "HUB-002", "HUB-003", "HUB-004", "HUB-005", "HUB-006", "HUB-007", "HUB-008", "HUB-009", "HUB-010"}

	composeAddress := func(c models.ContactInfo) string {
		return c.SubDistrict + ", " + c.District + ", " + c.Province
	}

	for i := 0; i < 20; i++ {
		orderNum := 10245 + i
		trackingNum := "TH2026" + padInt(100+i)

		cust := customers[i]
		recv := receivers[i]

		daysInMonth := func(m time.Month) int {
			switch m {
			case time.January, time.March, time.May, time.July, time.August, time.October, time.December:
				return 31
			case time.February:
				return 28
			default:
				return 30
			}
		}

		var month time.Month
		var shipmentStatus string
		var progress float64

		if i < 5 {
			// First 5 orders: force May with distinct statuses
			month = time.May
			mayStatuses := []string{"pending", "picked_up", "departed", "in_transit", "out_for_delivery"}
			shipmentStatus = mayStatuses[i]
			switch shipmentStatus {
			case "pending":
				progress = float64(rng.Intn(10))
			case "picked_up":
				progress = float64(rng.Intn(20) + 10)
			case "departed":
				progress = float64(rng.Intn(25) + 20)
			case "in_transit":
				progress = float64(rng.Intn(40) + 30)
			case "out_for_delivery":
				progress = float64(rng.Intn(25) + 65)
			}
		} else {
			// Remaining 15: random Jan-May
			month = time.Month(rng.Intn(5) + 1)
			if month < time.May {
				shipmentStatus = "delivered"
				progress = 100
			} else {
				roll := rng.Intn(100)
				switch {
				case roll < 10:
					shipmentStatus = "pending"
					progress = float64(rng.Intn(10))
				case roll < 20:
					shipmentStatus = "picked_up"
					progress = float64(rng.Intn(20) + 10)
				case roll < 30:
					shipmentStatus = "departed"
					progress = float64(rng.Intn(25) + 20)
				case roll < 55:
					shipmentStatus = "in_transit"
					progress = float64(rng.Intn(40) + 30)
				case roll < 75:
					shipmentStatus = "out_for_delivery"
					progress = float64(rng.Intn(25) + 65)
				case roll < 90:
					shipmentStatus = "delivered"
					progress = 100
				default:
					shipmentStatus = "delayed"
					progress = float64(rng.Intn(40) + 40)
				}
			}
		}

		day := rng.Intn(daysInMonth(month)) + 1
		createdAt := time.Date(2026, month, day, rng.Intn(12)+8, rng.Intn(60), 0, 0, time.UTC)

		hubID := hubs[rng.Intn(len(hubs))]
		if shipmentStatus == "pending" {
			hubID = ""
		}

		estDays := rng.Intn(5) + 1
		estDelivery := createdAt.AddDate(0, 0, estDays)
		carrier := "Thun-u-der Express"

		currLat := cust.Lat + (recv.Lat-cust.Lat)*progress/100.0
		currLng := cust.Lng + (recv.Lng-cust.Lng)*progress/100.0

		// Build events
		type event struct {
			timestamp time.Time
			status    string
			locName   string
			lat, lng  float64
			desc      string
		}
		var events []event
		maxDay := daysInMonth(month)

		makeEvent := func(dayOffset, hourMin, hourMax, minOffset int, status, locName string, lat, lng float64, desc string) {
			d := day + dayOffset
			if d > maxDay {
				d = maxDay
			}
			h := hourMin + rng.Intn(hourMax-hourMin+1)
			m := rng.Intn(60)
			t := time.Date(2026, month, d, h, m, 0, 0, time.UTC)
			events = append(events, event{t, status, locName, lat, lng, desc})
		}

		custAddr := cust.SubDistrict + ", " + cust.District + ", " + cust.Province
		recvAddr := recv.SubDistrict + ", " + recv.District + ", " + recv.Province

		// Label Created
		makeEvent(0, 7, 9, 0, "Label Created", custAddr, cust.Lat, cust.Lng, "Awaiting pickup.")

		if shipmentStatus != "pending" {
			// Picked Up
			makeEvent(rng.Intn(2), 9, 14, 0, "Picked Up", custAddr, cust.Lat, cust.Lng, "Parcel collected from sender.")
		}

		if shipmentStatus == "departed" || shipmentStatus == "in_transit" || shipmentStatus == "out_for_delivery" || shipmentStatus == "delivered" || shipmentStatus == "delayed" {
			// Departed
			makeEvent(rng.Intn(3)+1, 8, 13, 0, "Departed", "Bangkok Hub", 13.7563, 100.5018, "In transit to hub.")
		}

		if shipmentStatus == "in_transit" || shipmentStatus == "out_for_delivery" || shipmentStatus == "delivered" || shipmentStatus == "delayed" {
			// In Transit
			midLat := cust.Lat + (recv.Lat-cust.Lat)*0.5
			midLng := cust.Lng + (recv.Lng-cust.Lng)*0.5
			makeEvent(rng.Intn(3)+3, 8, 13, 0, "In Transit", "Central Hub", midLat, midLng, "Transit to next hub.")
		}

		if shipmentStatus == "out_for_delivery" || shipmentStatus == "delivered" || shipmentStatus == "delayed" {
			// Out for Delivery
			makeEvent(rng.Intn(2)+5, 8, 12, 0, "Out for Delivery", recvAddr, recv.Lat, recv.Lng, "Out for delivery.")
		}

		if shipmentStatus == "delivered" {
			// Delivered
			makeEvent(rng.Intn(2)+6, 9, 15, 0, "Delivered", recvAddr, recv.Lat, recv.Lng, "Delivered to recipient.")
		}

		if shipmentStatus == "delayed" {
			// Delayed
			makeEvent(rng.Intn(2)+6, 8, 12, 0, "Delayed", "Transit Hub", currLat, currLng, "Unexpected issue encountered.")
		}

		custContact := models.ContactInfo{
			Name: cust.Name, Zipcode: cust.Zipcode,
			SubDistrict: cust.SubDistrict, District: cust.District, Province: cust.Province,
			Coords: models.GeoPoint{Lat: cust.Lat, Lng: cust.Lng},
		}
		recvContact := models.ContactInfo{
			Name: recv.Name, Zipcode: recv.Zipcode,
			SubDistrict: recv.SubDistrict, District: recv.District, Province: recv.Province,
			Coords: models.GeoPoint{Lat: recv.Lat, Lng: recv.Lng},
		}

		shipment := models.Shipment{
			OrderID:           "ORD-" + padInt(orderNum),
			TrackingNumber:    trackingNum,
			Customer:          custContact,
			Receiver:          recvContact,
			Origin:            composeAddress(custContact),
			Destination:       composeAddress(recvContact),
			Status:            shipmentStatus,
			HubID:             hubID,
			Carrier:           carrier,
			Weight:            weights[i],
			Items:             rng.Intn(6) + 1,
			EstimatedDelivery: estDelivery,
			CreatedAt:         createdAt,
			Progress:          progress,
			CurrentCoords:     models.GeoPoint{Lat: currLat, Lng: currLng},
		}
		db.Create(&shipment)

		for _, e := range events {
			event := models.ShipmentEvent{
				ShipmentID:  shipment.ID,
				Status:      e.status,
				Location:    models.Location{Name: e.locName, Lat: e.lat, Lng: e.lng},
				Description: e.desc,
				CreatedAt:   e.timestamp,
			}
			db.Create(&event)
		}
	}
}

var weights = []float64{1.2, 2.5, 3.8, 5.2, 6.7, 8.5, 10.0, 4.2, 7.1, 3.3, 1.8, 5.5, 4.5, 2.1, 9.3, 6.0, 3.5, 7.8, 2.0, 4.8}

func padInt(n int) string {
	s := ""
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	for len(s) < 3 {
		s = "0" + s
	}
	return s
}
