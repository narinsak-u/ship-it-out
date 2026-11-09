package seed

import (
	"errors"
	"fmt"
	"strings"
	"time"

	authinfra "github.com/narinsak-u/backend/internal/auth/infrastructure"
	"gorm.io/gorm"
)

const adminCutoverKey = "legacy-admin-role-cutover-v1"

type bootstrapState struct {
	Key        string    `gorm:"primaryKey"`
	AdminEmail string    `gorm:"not null"`
	Completed  time.Time `gorm:"not null"`
}

func (bootstrapState) TableName() string { return "admin_bootstrap_states" }

const demoUserEmail = "demo@shipitout.example"
const demoUserPassword = "demo-read-only-2026"
const demoUserName = "Ship It Out Demo"

type demoBootstrapState struct {
	Key       string    `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
}

func (demoBootstrapState) TableName() string { return "demo_bootstrap_states" }

// BootstrapDemoUser provisions the public, read-only demo account once.
func BootstrapDemoUser(db *gorm.DB) error {
	if err := db.AutoMigrate(&demoBootstrapState{}); err != nil {
		return fmt.Errorf("migrate demo bootstrap state: %w", err)
	}
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT pg_advisory_xact_lock(?)", 130014).Error; err != nil {
			return fmt.Errorf("lock demo bootstrap: %w", err)
		}
		var state demoBootstrapState
		err := tx.First(&state, "key = ?", "public-demo-user-v1").Error
		if err == nil {
			var user authinfra.UserRecord
			if err := tx.First(&user, state.UserID).Error; err != nil || user.Email != demoUserEmail || user.Role != "customer" {
				return errors.New("demo bootstrap is complete but demo account is unavailable")
			}
			if !authinfra.CheckPassword(demoUserPassword, user.Password) {
				return errors.New("demo bootstrap account credentials changed")
			}
			return nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("read demo bootstrap state: %w", err)
		}
		var collision authinfra.UserRecord
		if err := tx.Where("email = ?", demoUserEmail).First(&collision).Error; err == nil {
			return errors.New("demo email already belongs to an account")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("check demo email collision: %w", err)
		}
		hash, err := authinfra.HashPassword(demoUserPassword)
		if err != nil {
			return fmt.Errorf("hash demo password: %w", err)
		}
		user := authinfra.UserRecord{Name: demoUserName, Email: demoUserEmail, Password: hash, Role: "customer"}
		if err := tx.Create(&user).Error; err != nil {
			return fmt.Errorf("create demo account: %w", err)
		}
		if err := tx.Create(&demoBootstrapState{Key: "public-demo-user-v1", UserID: user.ID, CreatedAt: time.Now().UTC()}).Error; err != nil {
			return fmt.Errorf("record demo bootstrap: %w", err)
		}
		return nil
	})
}

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
			var admin authinfra.UserRecord
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

		var collision authinfra.UserRecord
		if err := tx.Where("email = ?", email).First(&collision).Error; err == nil {
			return fmt.Errorf("admin email already belongs to an account")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("check admin email collision: %w", err)
		}
		hash, err := authinfra.HashPassword(password)
		if err != nil {
			return fmt.Errorf("hash admin password: %w", err)
		}
		if err := tx.Model(&authinfra.UserRecord{}).Where("role = ?", "admin").Update("role", "customer").Error; err != nil {
			return fmt.Errorf("demote legacy admins: %w", err)
		}
		admin := authinfra.UserRecord{Name: name, Email: email, Password: hash, Role: "admin"}
		if err := tx.Create(&admin).Error; err != nil {
			return fmt.Errorf("create configured admin: %w", err)
		}
		if err := tx.Create(&bootstrapState{Key: adminCutoverKey, AdminEmail: email, Completed: time.Now().UTC()}).Error; err != nil {
			return fmt.Errorf("record admin cutover: %w", err)
		}
		return nil
	})
}
