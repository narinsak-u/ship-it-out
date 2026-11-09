package infrastructure

import "time"

type UserRecord struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Email     string    `gorm:"unique;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	Role      string    `gorm:"not null;default:customer" json:"role"`
	CreatedAt time.Time `json:"created_at"`
}
func (UserRecord) TableName() string { return "users" }
