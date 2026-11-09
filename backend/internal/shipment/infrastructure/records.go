package infrastructure

import "time"


type shipmentRecord struct {
	ID uint `gorm:"primaryKey;autoIncrement"`
	OrderID string `gorm:"uniqueIndex;not null"`
	TrackingNumber string `gorm:"unique;index;not null"`
	CustomerName string `gorm:"column:customer_name"`
	CustomerZipcode string `gorm:"column:customer_zipcode"`
	CustomerSubDistrict string `gorm:"column:customer_sub_district"`
	CustomerDistrict string `gorm:"column:customer_district"`
	CustomerProvince string `gorm:"column:customer_province"`
	CustomerLat float64 `gorm:"column:customer_lat"`
	CustomerLng float64 `gorm:"column:customer_lng"`
	ReceiverName string `gorm:"column:receiver_name"`
	ReceiverZipcode string `gorm:"column:receiver_zipcode"`
	ReceiverSubDistrict string `gorm:"column:receiver_sub_district"`
	ReceiverDistrict string `gorm:"column:receiver_district"`
	ReceiverProvince string `gorm:"column:receiver_province"`
	ReceiverLat float64 `gorm:"column:receiver_lat"`
	ReceiverLng float64 `gorm:"column:receiver_lng"`
	CurrentLat float64 `gorm:"column:current_lat"`
	CurrentLng float64 `gorm:"column:current_lng"`
	Origin string `gorm:"not null"`
	Destination string `gorm:"not null"`
	Status string `gorm:"index;not null;default:pending"`
	Carrier string
	HubID string `gorm:"column:hub_id"`
	Weight float64
	Items int
	EstimatedDelivery time.Time
	CreatedAt time.Time `gorm:"index"`
	Progress float64
}
func (shipmentRecord) TableName() string { return "shipments" }

type eventRecord struct {
	ID uint `gorm:"primaryKey;autoIncrement"`
	ShipmentID uint `gorm:"not null;index"`
	Status string `gorm:"not null"`
	LocationName string `gorm:"column:location_name"`
	LocationLat float64 `gorm:"column:location_lat"`
	LocationLng float64 `gorm:"column:location_lng"`
	Description string
	CreatedAt time.Time
}
func (eventRecord) TableName() string { return "shipment_events" }
