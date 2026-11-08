package seed

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/narinsak-u/backend/internal/models"
	"github.com/narinsak-u/backend/pkg/utils"
	"gorm.io/gorm"
)

const adminCutoverKey = "legacy-admin-role-cutover-v1"

type bootstrapState struct {
	Key        string    `gorm:"primaryKey"`
	AdminEmail string    `gorm:"not null"`
	Completed  time.Time `gorm:"not null"`
}

func (bootstrapState) TableName() string { return "admin_bootstrap_states" }

// BootstrapAdmin demotes legacy admins and creates the configured admin exactly once.
func BootstrapAdmin(db *gorm.DB, name, email, password string, cutoverEnabled bool) error {
	name, email = strings.TrimSpace(name), strings.TrimSpace(email)
	if name == "" || email == "" || password == "" {
		return errors.New("admin bootstrap requires name, email, and password")
	}
	if err := db.AutoMigrate(&bootstrapState{}); err != nil {
		return fmt.Errorf("migrate admin cutover state: %w", err)
	}
	return db.Transaction(func(tx *gorm.DB) error {
		// Serialize concurrent application starts before checking the durable marker.
		if err := tx.Exec("SELECT pg_advisory_xact_lock(?)", 130013).Error; err != nil {
			return fmt.Errorf("lock admin cutover: %w", err)
		}
		var state bootstrapState
		err := tx.First(&state, "key = ?", adminCutoverKey).Error
		if err == nil {
			if state.AdminEmail != email {
				return errors.New("configured admin email differs from completed cutover")
			}
			var admin models.User
			if err := tx.Where("email = ? AND role = ?", state.AdminEmail, "admin").First(&admin).Error; err != nil {
				return fmt.Errorf("cutover is complete but configured admin is unavailable: %w", err)
			}
			return nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("read admin cutover state: %w", err)
		}
		if !cutoverEnabled {
			return errors.New("admin role cutover is incomplete; enable ADMIN_ROLE_CUTOVER after a verified backup")
		}

		var collision models.User
		if err := tx.Where("email = ?", email).First(&collision).Error; err == nil {
			return fmt.Errorf("admin email already belongs to an account")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("check admin email collision: %w", err)
		}
		hash, err := utils.HashPassword(password)
		if err != nil {
			return fmt.Errorf("hash admin password: %w", err)
		}
		if err := tx.Model(&models.User{}).Where("role = ?", "admin").Update("role", "customer").Error; err != nil {
			return fmt.Errorf("demote legacy admins: %w", err)
		}
		admin := models.User{Name: name, Email: email, Password: hash, Role: "admin"}
		if err := tx.Create(&admin).Error; err != nil {
			return fmt.Errorf("create configured admin: %w", err)
		}
		if err := tx.Create(&bootstrapState{Key: adminCutoverKey, AdminEmail: email, Completed: time.Now().UTC()}).Error; err != nil {
			return fmt.Errorf("record admin cutover: %w", err)
		}
		return nil
	})
}
