package infrastructure

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("myPassword123")
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, "myPassword123", hash)
}

func TestCheckPassword(t *testing.T) {
	hash, err := HashPassword("correctPassword")
	require.NoError(t, err)
	assert.True(t, CheckPassword("correctPassword", hash))
	assert.False(t, CheckPassword("wrongPassword", hash))
}

func TestCheckPassword_WrongHash(t *testing.T) {
	assert.False(t, CheckPassword("any", "invalidhash"))
}
