package auth

import (
	"errors"
	"testing"
	"time"

	"github.com/narinsak-u/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type serviceRepo struct {
	user    *models.User
	create  error
	findErr error
}

func (r *serviceRepo) Create(user *models.User) error {
	r.user = user
	if r.create == nil {
		user.ID = 4
	}
	return r.create
}

func (r *serviceRepo) FindByEmail(string) (*models.User, error) { return r.user, r.findErr }
func (r *serviceRepo) FindByID(uint) (*models.User, error)       { return r.user, r.findErr }

func TestServiceRegisterForcesCustomerAndHashesPassword(t *testing.T) {
	repo := &serviceRepo{}
	service := NewService(repo, "secret", time.Hour)

	user, token, err := service.Register("N", "n@example.com", "password")
	require.NoError(t, err)
	assert.Equal(t, "customer", user.Role)
	assert.NotEqual(t, "password", repo.user.Password)
	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(repo.user.Password), []byte("password")))
	assert.NotEmpty(t, token.String())
	claims, err := ParseToken(token.String(), []byte("secret"))
	require.NoError(t, err)
	assert.Equal(t, uint(4), claims.UserID)
	assert.Equal(t, "customer", claims.Role)
}

func TestServiceLoginRejectsInvalidCredentialsAndPropagatesRepositoryErrors(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.MinCost)
	require.NoError(t, err)
	repo := &serviceRepo{user: &models.User{ID: 4, Password: string(hash)}}
	service := NewService(repo, "secret", time.Hour)
	user, token, err := service.Login("n@example.com", "correct")
	require.NoError(t, err)
	assert.Equal(t, uint(4), user.ID)
	assert.NotEmpty(t, token.String())

	_, _, err = service.Login("n@example.com", "wrong")
	assert.ErrorIs(t, err, ErrInvalidCredentials)

	repo.findErr = errors.New("database unavailable")
	_, _, err = service.Login("n@example.com", "correct")
	assert.ErrorIs(t, err, repo.findErr)
}

func TestServiceRegisterPropagatesRepositoryErrors(t *testing.T) {
	repo := &serviceRepo{create: errors.New("database unavailable")}
	_, _, err := NewService(repo, "secret", time.Hour).Register("N", "n@example.com", "password")
	assert.ErrorIs(t, err, repo.create)
}
