package analytics

import (
	"github.com/gofiber/fiber/v2"
	"github.com/narinsak-u/backend/internal/database"
	"github.com/narinsak-u/backend/internal/models"
	"github.com/narinsak-u/backend/pkg/utils"
)

func Overview(c *fiber.Ctx) error {
	var total int64
	database.DB.Model(&models.Shipment{}).Count(&total)

	var active int64
	database.DB.Model(&models.Shipment{}).Where("status NOT IN ?", []string{"DELIVERED", "RETURNED"}).Count(&active)

	var delivered int64
	database.DB.Model(&models.Shipment{}).Where("status = ?", "DELIVERED").Count(&delivered)

	type StatusCount struct {
		Status string
		Count  int64
	}
	var byStatus []StatusCount
	database.DB.Model(&models.Shipment{}).Select("status, count(*) as count").Group("status").Scan(&byStatus)

	return utils.Success(c, fiber.Map{
		"total":     total,
		"active":    active,
		"delivered": delivered,
		"by_status": byStatus,
	})
}
