package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zjr71163356/simplebank/utils"
)

func CreateAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    utils.RandomString(),
		Balance:  123,
		Currency: "USD",
	}
	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return account
}

func GetAccount(t *testing.T, account Account) (Account, error) {

	account2, err := testQueries.GetAccount(context.Background(), account.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account.Owner, account2.Owner)
	require.Equal(t, account.Currency, account2.Currency)
	require.Equal(t, account.Balance, account2.Balance)
	require.Equal(t, account.ID, account2.ID)
	require.Equal(t, account.CreatedAt, account2.CreatedAt)

	return account2, err
}
func TestCreateAccount(t *testing.T) {
	CreateAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := CreateAccount(t)
	GetAccount(t, account1)

}

func TestDeleteAccount(t *testing.T) {

	account := CreateAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account.ID)

	require.NoError(t, err)

}

func TestListAccounts(t *testing.T) {

	utils.RandomString(5)
	// Create 10 random accounts
	var createdAccounts []Account

	for i := 0; i < 10; i++ {
		arg := CreateAccountParams{
			Owner:    utils.RandomString(6),
			Balance:  utils.RandomInt(0, 1000),
			Currency: utils.RandomCurrency(),
		}
		account, err := testQueries.CreateAccount(context.Background(), arg)
		require.NoError(t, err)
		createdAccounts = append(createdAccounts, account)
	}

	// Test listing accounts
	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
