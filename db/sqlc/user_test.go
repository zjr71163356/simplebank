package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zjr71163356/simplebank/utils"
)

func createRandomUser(t *testing.T) (User, error) {
	arg := CreateUserParams{
		Username:       utils.RandomOwnerName(),
		HashedPassword: "secret",
		FullName:       utils.RandomString(5),
		Email:          utils.RandomEmail(),
	}
	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)

	require.Zero(t, user.PasswordChangedAt)
	require.NotEmpty(t, user.CreatedAt)

	return user, err
}
func createUserWithFixedName(t *testing.T, username string) (User, error) {
	arg := CreateUserParams{
		Username:       username,
		HashedPassword: "secret",
		FullName:       utils.RandomString(5),
		Email:          utils.RandomEmail(),
	}
	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)

	require.Zero(t, user.PasswordChangedAt)
	require.NotEmpty(t, user.CreatedAt)

	return user, err
}
func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	createdUser, err := createRandomUser(t)

	gotUser, err := testQueries.GetUser(context.Background(), createdUser.Username)
	require.NoError(t, err)
	require.NotEmpty(t, gotUser)

	require.Equal(t, createdUser.Username, gotUser.Username)
	require.Equal(t, createdUser.Email, gotUser.Email)
	require.Equal(t, createdUser.FullName, gotUser.FullName)
	require.Equal(t, createdUser.HashedPassword, gotUser.HashedPassword)
	require.WithinDuration(t, createdUser.PasswordChangedAt, gotUser.PasswordChangedAt, time.Second)
	require.WithinDuration(t, createdUser.CreatedAt, gotUser.CreatedAt, time.Second)

}
