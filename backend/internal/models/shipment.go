package models

import "time"

type Shipment struct {
	ID                 uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	TrackingNumber     string    `gorm:"unique;index;not null" json:"tracking_number"`
	SenderName         string    `gorm:"not null" json:"sender_name"`
	ReceiverName       string    `gorm:"not null" json:"receiver_name"`
	OriginAddress      string    `gorm:"not null" json:"origin_address"`
	DestinationAddress string    `gorm:"not null" json:"destination_address"`
	Weight             float64   `gorm:"not null" json:"weight"`
	Status             string    `gorm:"not null;default:CREATED" json:"status"`
	EstimatedDelivery  time.Time `json:"estimated_delivery"`
	CreatedAt          time.Time `json:"created_at"`
}
