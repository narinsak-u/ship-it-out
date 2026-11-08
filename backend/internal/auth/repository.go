// Package auth provides user registration, login, session verification, and logout.
package auth

import "github.com/narinsak-u/backend/internal/models"

// Repository defines the data-access contract for user accounts. The auth service depends on
// this interface, so tests can swap in a fake without connecting to a real database.
type Repository interface {
	// Create inserts a new user. Returns ErrEmailAlreadyRegistered for a duplicate email.
	Create(user *models.User) error
	// FindByEmail looks up a user by email. Returns ErrUserNotFound if absent.
	FindByEmail(email string) (*models.User, error)
	// FindByID looks up a user by primary key. Returns ErrUserNotFound if absent.
	FindByID(id uint) (*models.User, error)
}
