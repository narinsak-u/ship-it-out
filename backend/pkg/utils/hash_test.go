package utils_test

import (
	"testing"

	"github.com/narinsak-u/backend/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	hash, err := utils.HashPassword("myPassword123")
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, "myPassword123", hash)
}

func TestCheckPassword(t *testing.T) {
	hash, err := utils.HashPassword("correctPassword")
	require.NoError(t, err)
	assert.True(t, utils.CheckPassword("correctPassword", hash))
	assert.False(t, utils.CheckPassword("wrongPassword", hash))
}

func TestCheckPassword_WrongHash(t *testing.T) {
	assert.False(t, utils.CheckPassword("any", "invalidhash"))
}
