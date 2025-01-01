package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("password")

	require.NoError(t, err, "error hashing password")
	require.NotEmpty(t, hash, "hash cannot be empty")

	assert.NotEqual(t, "password", hash, "inputed pass and hash should not be equal")

}

func TestComparePasswords(t *testing.T) {
	hash, err := HashPassword("password")

	require.NoError(t, err, "error comparing passwords")

	assert.True(t, ComparePasswords(hash, []byte("password")), "expected password to match hash")
	assert.False(t, ComparePasswords(hash, []byte("notpassword")), "expected password to not match hash")
}
