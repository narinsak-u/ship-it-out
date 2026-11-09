package domain

import "time"

const (
	RoleCustomer = "customer"
	RoleStaff    = "staff"
	RoleAdmin    = "admin"
)

type User struct {
	ID        uint
	Name      string
	Email     string
	Password  string
	Role      string
	CreatedAt time.Time
}
