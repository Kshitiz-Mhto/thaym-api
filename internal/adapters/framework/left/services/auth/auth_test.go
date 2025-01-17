package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateJWT(t *testing.T) {

	secret := []byte("secret")
	token, err := CreateJWT(secret, "abc")

	require.NoError(t, err, "error creating JWT")

	assert.NotEmpty(t, token, "expected token to be not empty")
}
