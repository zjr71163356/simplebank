package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zjr71163356/simplebank/utils"
)

func TestJWTMaker(t *testing.T) {
	secretKey := utils.RandomString(32)
	jwtMaker, err := NewJWTMaker(secretKey)
	require.NoError(t, err)

	username := utils.RandomOwnerName()
	duration := time.Minute

	IssuedAt := time.Now()
	ExpiredAt := time.Now().Add(duration)

	token, err := jwtMaker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := jwtMaker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotEmpty(t, payload.Id)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, IssuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, ExpiredAt, payload.ExpiredAt, time.Second)

}
