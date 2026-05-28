package analytics

import (
	"github.com/gofiber/fiber/v2"
	"github.com/narinsak-u/backend/internal/database"
	"github.com/narinsak-u/backend/internal/models"
	"github.com/narinsak-u/backend/pkg/utils"
)

var provinceRegion = map[string]string{
	"Bangkok": "Central", "Nonthaburi": "Central", "Pathum Thani": "Central",
	"Samut Prakan": "Central", "Phra Nakhon Si Ayutthaya": "Central",
	"Nakhon Pathom": "Central", "Samut Sakhon": "Central", "Samut Songkhram": "Central",
	"Sing Buri": "Central", "Ang Thong": "Central", "Lopburi": "Central",
	"Saraburi": "Central", "Chai Nat": "Central", "Suphan Buri": "Central",
	"Nakhon Sawan": "Central", "Uthai Thani": "Central", "Phichit": "Central",
	"Chonburi": "East", "Rayong": "East", "Chachoengsao": "East",
	"Chanthaburi": "East", "Trat": "East", "Prachinburi": "East",
	"Sa Kaeo": "East", "Nakhon Nayok": "East",
	"Chiang Mai": "North", "Chiang Rai": "North", "Lampang": "North",
	"Lamphun": "North", "Phrae": "North", "Nan": "North", "Phayao": "North",
	"Mae Hong Son": "North", "Uttaradit": "North",
	"Kamphaeng Phet": "North", "Phetchabun": "North", "Phitsanulok": "North", "Sukhothai": "North",
	"Kanchanaburi": "West", "Phetchaburi": "West", "Prachuap Khiri Khan": "West",
	"Ratchaburi": "West", "Tak": "West",
	"Khon Kaen": "North-east", "Nakhon Ratchasima": "North-east",
	"Udon Thani": "North-east", "Ubon Ratchathani": "North-east",
	"Buriram": "North-east", "Surin": "North-east", "Sisaket": "North-east",
	"Chaiyaphum": "North-east", "Loei": "North-east", "Nong Bua Lamphu": "North-east",
	"Nong Khai": "North-east", "Maha Sarakham": "North-east", "Roi Et": "North-east",
	"Kalasin": "North-east", "Sakon Nakhon": "North-east", "Nakhon Phanom": "North-east",
	"Mukdahan": "North-east", "Amnat Charoen": "North-east", "Bueng Kan": "North-east",
	"Yasothon": "North-east",
	"Phuket":   "South", "Songkhla": "South", "Surat Thani": "South",
	"Nakhon Si Thammarat": "South", "Krabi": "South", "Trang": "South",
	"Phatthalung": "South", "Satun": "South", "Chumphon": "South",
	"Ranong": "South", "Phangnga": "South", "Pattani": "South",
	"Yala": "South", "Narathiwat": "South",
}

type provinceCount struct {
	Province  string
	Total     int64
	Delivered int64
}

type regionCount struct {
	Name      string `json:"name"`
	Total     int64  `json:"total"`
	Delivered int64  `json:"delivered"`
}

func Overview(c *fiber.Ctx) error {
	var total int64
	database.DB.Model(&models.Shipment{}).Count(&total)

	var active int64
	database.DB.Model(&models.Shipment{}).Where("status NOT IN ?", []string{"DELIVERED", "RETURNED"}).Count(&active)

	var delivered int64
	database.DB.Model(&models.Shipment{}).Where("status = ?", "DELIVERED").Count(&delivered)

	type StatusCount struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}
	var byStatus []StatusCount
	database.DB.Model(&models.Shipment{}).Select("status, count(*) as count").Group("status").Scan(&byStatus)

	var byProvince []provinceCount
	database.DB.Model(&models.Shipment{}).
		Select("receiver_province as province, count(*) as total, sum(case when status = 'DELIVERED' then 1 else 0 end) as delivered").
		Group("receiver_province").
		Scan(&byProvince)

	regionMap := make(map[string]*regionCount)
	for _, p := range byProvince {
		region := provinceRegion[p.Province]
		if region == "" {
			region = "Other"
		}
		r, ok := regionMap[region]
		if !ok {
			r = &regionCount{Name: region}
			regionMap[region] = r
		}
		r.Total += p.Total
		r.Delivered += p.Delivered
	}

	byRegion := make([]regionCount, 0, len(regionMap))
	for _, r := range regionMap {
		byRegion = append(byRegion, *r)
	}

	return utils.Success(c, fiber.Map{
		"total":     total,
		"active":    active,
		"delivered": delivered,
		"by_status": byStatus,
		"by_region": byRegion,
	})
}
