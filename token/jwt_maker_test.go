package token

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
	ExpiredAt := IssuedAt.Add(duration)

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

func TestExpiredJWTToken(t *testing.T) {
	secretKey := utils.RandomString(32)
	jwtMaker, err := NewJWTMaker(secretKey)
	require.NoError(t, err)

	username := utils.RandomOwnerName()
	duration := time.Minute

	token, err := jwtMaker.CreateToken(username, -duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := jwtMaker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, jwt.ErrTokenExpired.Error())
	require.Nil(t, payload)

}

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	secretKey := utils.RandomString(32)
	username := utils.RandomOwnerName()
	duration := time.Minute

	payload, err := NewPayload(username, duration)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	jwtMaker, err := NewJWTMaker(secretKey)
	require.NoError(t, err)
	payload, err = jwtMaker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
