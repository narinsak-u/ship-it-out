package shipment

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/narinsak-u/backend/internal/database"
	"github.com/narinsak-u/backend/internal/models"
	"github.com/narinsak-u/backend/pkg/utils"
)

type CreateRequest struct {
	SenderName         string  `json:"sender_name"`
	ReceiverName       string  `json:"receiver_name"`
	OriginAddress      string  `json:"origin_address"`
	DestinationAddress string  `json:"destination_address"`
	Weight             float64 `json:"weight"`
}

func generateTrackingNumber() string {
	return fmt.Sprintf("TH%d%05d", time.Now().Year(), time.Now().UnixMilli()%100000)
}

func List(c *fiber.Ctx) error {
	var shipments []models.Shipment
	database.DB.Find(&shipments)
	return utils.Success(c, shipments)
}

func Create(c *fiber.Ctx) error {
	var req CreateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, 400, "invalid request body")
	}

	shipment := models.Shipment{
		TrackingNumber:     generateTrackingNumber(),
		SenderName:         req.SenderName,
		ReceiverName:       req.ReceiverName,
		OriginAddress:      req.OriginAddress,
		DestinationAddress: req.DestinationAddress,
		Weight:             req.Weight,
		Status:             "CREATED",
	}

	if result := database.DB.Create(&shipment); result.Error != nil {
		return utils.Error(c, 500, "failed to create shipment")
	}

	return utils.Success(c, shipment)
}

func GetByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return utils.Error(c, 400, "invalid id")
	}

	var shipment models.Shipment
	if result := database.DB.First(&shipment, id); result.Error != nil {
		return utils.Error(c, 404, "shipment not found")
	}

	return utils.Success(c, shipment)
}

func UpdateStatus(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return utils.Error(c, 400, "invalid id")
	}

	var body struct {
		Status string `json:"status"`
	}
	if err := c.BodyParser(&body); err != nil {
		return utils.Error(c, 400, "invalid request body")
	}

	var shipment models.Shipment
	if result := database.DB.First(&shipment, id); result.Error != nil {
		return utils.Error(c, 404, "shipment not found")
	}

	shipment.Status = body.Status
	database.DB.Save(&shipment)

	event := models.ShipmentEvent{
		ShipmentID:  shipment.ID,
		Status:      body.Status,
		Description: fmt.Sprintf("Status updated to %s", body.Status),
	}
	database.DB.Create(&event)

	return utils.Success(c, shipment)
}
