package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := RandomString(6)

	hashedPassword1, err := HashPassWord(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword1)

	err = MatchPassWord(hashedPassword1, password)
	require.NoError(t, err)

	wrongPassword := RandomString(6)
	err = MatchPassWord(hashedPassword1, wrongPassword)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	hashedPassword2, err := HashPassWord(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword2)
	require.NotEqual(t, hashedPassword1, hashedPassword2)
}
