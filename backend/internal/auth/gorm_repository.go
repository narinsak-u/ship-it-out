package auth

import (
	"github.com/narinsak-u/backend/internal/models"
	"gorm.io/gorm"
)

// GormRepository implements Repository using GORM as the backend. It talks directly to Postgres.
// Create one via NewGormRepository — the zero value is not usable.
type GormRepository struct{ db *gorm.DB }

// NewGormRepository opens a GORM-backed auth repository using the given *gorm.DB connection.
func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

// Create inserts a new user row into Postgres. Returns error (typically a unique-constraint
// violation) if a user with the same email already exists. The password field must already be
// hashed by the caller.
func (r *GormRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// FindByEmail looks up a user by their email address. Returns nil + error if not found.
// The caller should NOT distinguish "user not found" from "wrong password" in error messages
// (see handler Login — both return 401 "invalid email or password").
func (r *GormRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID looks up a user by their auto-increment primary key. Used by the Me endpoint
// after the AuthRequired middleware extracts user_id from the JWT claims.
func (r *GormRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
