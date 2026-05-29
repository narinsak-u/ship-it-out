// Package analytics provides HTTP handlers for dashboard aggregate statistics.
// Overview returns total, active, delivered counts, status distribution, and
// geographic breakdown by Thai region. TimeSeries returns monthly and day-of-week
// trends used to power the frontend analytics charts.
package analytics

import (
	"github.com/gofiber/fiber/v2"
	"github.com/narinsak-u/backend/internal/data"
	"github.com/narinsak-u/backend/internal/database"
	"github.com/narinsak-u/backend/internal/models"
	"github.com/narinsak-u/backend/pkg/utils"
)

type provinceCount struct {
	Province string
	Total    int64
}

type regionCount struct {
	Name  string `json:"name"`
	Total int64  `json:"total"`
}

func Overview(c *fiber.Ctx) error {
	var total int64
	database.DB.Model(&models.Shipment{}).Count(&total)

	var active int64
	database.DB.Model(&models.Shipment{}).Where("status NOT IN ?", []string{"delivered", "returned"}).Count(&active)

	var delivered int64
	database.DB.Model(&models.Shipment{}).Where("status = ?", "delivered").Count(&delivered)

	type StatusCount struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}
	var byStatus []StatusCount
	database.DB.Model(&models.Shipment{}).Select("status, count(*) as count").Group("status").Scan(&byStatus)

	var byProvince []provinceCount
	database.DB.Model(&models.Shipment{}).
		Select("receiver_province as province, count(*) as total").
		Group("receiver_province").
		Scan(&byProvince)

	regionMap := map[string]*regionCount{
		"Central":    {Name: "Central"},
		"East":       {Name: "East"},
		"North":      {Name: "North"},
		"West":       {Name: "West"},
		"North-east": {Name: "North-east"},
		"South":      {Name: "South"},
	}
	for _, p := range byProvince {
		region := data.ThailandProvinceRegion[p.Province]
		if region == "" {
			continue
		}
		r := regionMap[region]
		r.Total += p.Total
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

type monthCount struct {
	Month string `json:"month"`
	Count int64  `json:"count"`
}

type dayCount struct {
	Day   string `json:"day"`
	Count int64  `json:"count"`
}

func TimeSeries(c *fiber.Ctx) error {
	var byMonth []monthCount
	database.DB.Model(&models.Shipment{}).
		Select("to_char(created_at, 'YYYY-MM') as month, count(*) as count").
		Group("month").
		Order("month").
		Scan(&byMonth)

	var byDay []dayCount
	database.DB.Model(&models.Shipment{}).
		Select("trim(to_char(created_at, 'Day')) as day, count(*) as count").
		Group("day").
		Order("min(extract(dow from created_at))").
		Scan(&byDay)

	return utils.Success(c, fiber.Map{
		"by_month":       byMonth,
		"by_day_of_week": byDay,
	})
}
