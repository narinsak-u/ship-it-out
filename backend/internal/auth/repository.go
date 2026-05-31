// Package auth provides user registration, login, session verification, and logout.
package auth

import "github.com/narinsak-u/backend/internal/models"

// Repository defines the data-access contract for user accounts. The auth Handler depends on
// this interface, so tests can swap in a mock without connecting to a real database.
type Repository interface {
	// Create inserts a new user. Returns error on duplicate email (unique constraint).
	Create(user *models.User) error
	// FindByEmail looks up a user by email (for login). Returns error if not found.
	FindByEmail(email string) (*models.User, error)
	// FindByID looks up a user by primary key (for session/profile retrieval).
	FindByID(id uint) (*models.User, error)
}
