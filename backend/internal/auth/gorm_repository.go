package auth

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
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
	if err := r.db.Create(user).Error; err != nil {
		if isDuplicateEmail(err) {
			return ErrEmailAlreadyRegistered
		}
		return err
	}
	return nil
}

// FindByEmail looks up a user by their email address. Returns ErrUserNotFound if absent.
func (r *GormRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// FindByID looks up a user by their auto-increment primary key. Used by the Me endpoint
// after the AuthRequired middleware extracts user_id from the JWT claims.
func (r *GormRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func isDuplicateEmail(err error) bool {
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
