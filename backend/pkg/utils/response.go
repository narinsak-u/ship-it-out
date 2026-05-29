package utils

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/narinsak-u/backend/internal/models"
)

// ComposeAddress builds a display string from a ContactInfo's address fields.
func ComposeAddress(c models.ContactInfo) string {
	return fmt.Sprintf("%s, %s, %s", c.SubDistrict, c.District, c.Province)
}

func Success(c *fiber.Ctx, data interface{}) error {
	return c.JSON(fiber.Map{"success": true, "data": data})
}

func Error(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{"success": false, "error": message})
}

func SuccessWithPagination(c *fiber.Ctx, data interface{}, page, limit, total int) error {
	totalPages := (total + limit - 1) / limit
	if totalPages < 1 {
		totalPages = 1
	}
	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
		"pagination": fiber.Map{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": totalPages,
		},
	})
}
