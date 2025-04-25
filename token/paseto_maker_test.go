package token

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
	"github.com/zjr71163356/simplebank/utils"
)

func TestPasetoMaker(t *testing.T) {
	secretKey := utils.RandomString(32)
	maker, err := NewPasetoMaker(secretKey)
	require.NoError(t, err)

	username := utils.RandomOwnerName()
	duration := time.Minute

	IssuedAt := time.Now()
	ExpiredAt := IssuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotEmpty(t, payload.Id)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, IssuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, ExpiredAt, payload.ExpiredAt, time.Second)

}

func TestExpiredPasetoToken(t *testing.T) {
	secretKey := utils.RandomString(32)
	maker, err := NewPasetoMaker(secretKey)
	require.NoError(t, err)

	username := utils.RandomOwnerName()
	duration := time.Minute

	token, err := maker.CreateToken(username, -duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, jwt.ErrTokenExpired.Error())
	require.Nil(t, payload)

}
