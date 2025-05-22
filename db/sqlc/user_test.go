package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zjr71163356/simplebank/utils"
)

func createRandomUser(t *testing.T) (User, error) {

	HashedPassword, err := utils.RandomHashPassWord()
	arg := CreateUserParams{
		Username:       utils.RandomOwnerName(),
		HashedPassword: HashedPassword,
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
	require.NoError(t, err)
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

func TestUpdateUserOnlyFullName(t *testing.T) {
	oldUser, err := createRandomUser(t)
	require.NoError(t, err)

	newFullName := utils.RandomOwnerName()
	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
	})
	require.NoError(t, err)

	require.NotEqual(t, oldUser.FullName, updatedUser.FullName)
	require.Equal(t, newFullName, updatedUser.FullName)
	require.Equal(t, oldUser.Email, updatedUser.Email)
	require.Equal(t, oldUser.HashedPassword, updatedUser.HashedPassword)

}

func TestUpdateUserOnlyEmail(t *testing.T) {
	oldUser, err := createRandomUser(t)
	require.NoError(t, err)

	newEmail := utils.RandomEmail()
	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
	})
	require.NoError(t, err)

	require.NotEqual(t, oldUser.Email, updatedUser.Email)
	require.Equal(t, newEmail, updatedUser.Email)
	require.Equal(t, oldUser.FullName, updatedUser.FullName)
	require.Equal(t, oldUser.HashedPassword, updatedUser.HashedPassword)

}

func TestUpdateUserOnlyPassword(t *testing.T) {
	oldUser, err := createRandomUser(t)
	require.NoError(t, err)

	newPassword := utils.RandomString(10)
	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: sql.NullString{
			String: newPassword,
			Valid:  true,
		},
	})
	require.NoError(t, err)

	require.NotEqual(t, oldUser.HashedPassword, updatedUser.HashedPassword)
	require.Equal(t, newPassword, updatedUser.HashedPassword)
	require.Equal(t, oldUser.FullName, updatedUser.FullName)
	require.Equal(t, oldUser.Email, updatedUser.Email)

}

func TestUpdateUserAllFields(t *testing.T) {
	oldUser, err := createRandomUser(t)
	require.NoError(t, err)

	newFullName := utils.RandomOwnerName()
	newEmail := utils.RandomEmail()
	newPassword := utils.RandomString(10)

	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
		HashedPassword: sql.NullString{
			String: newPassword,
			Valid:  true,
		},
	})
	require.NoError(t, err)

	require.NotEqual(t, oldUser.FullName, updatedUser.FullName)
	require.Equal(t, newFullName, updatedUser.FullName)
	require.NotEqual(t, oldUser.Email, updatedUser.Email)
	require.Equal(t, newEmail, updatedUser.Email)
	require.NotEqual(t, oldUser.HashedPassword, updatedUser.HashedPassword)
	require.Equal(t, newPassword, updatedUser.HashedPassword)

}
