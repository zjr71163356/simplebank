package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zjr71163356/simplebank/utils"
)

func checkCreatedAccount(t *testing.T, account Account, arg CreateAccountParams) {
	require.NotEmpty(t, account)

	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
}

func checkGotAccount(t *testing.T, account Account, account2 Account) {
	require.NotEmpty(t, account2)

	require.Equal(t, account.Owner, account2.Owner)
	require.Equal(t, account.Currency, account2.Currency)
	require.Equal(t, account.Balance, account2.Balance)
	require.Equal(t, account.ID, account2.ID)
	require.Equal(t, account.CreatedAt, account2.CreatedAt)
}

func createRandomAccount(t *testing.T) (Account, error) {
	user, err := createRandomUser(t)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  utils.RandomInt63(0, 5000),
		Currency: utils.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	checkCreatedAccount(t, account, arg)
	return account, err
}

func CreateAccountWithFixedOwner(t *testing.T, owner string) Account {
	arg := CreateAccountParams{
		Owner:    owner,
		Balance:  utils.RandomInt63(0, 5000),
		Currency: utils.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	checkCreatedAccount(t, account, arg)
	return account
}

func GetAccount(t *testing.T, account Account) (Account, error) {

	account2, err := testQueries.GetAccount(context.Background(), account.ID)
	require.NoError(t, err)
	checkGotAccount(t, account, account2)
	return account2, err
}

func GetAccountForUpdate(t *testing.T, account Account) (Account, error) {
	account2, err := testQueries.GetAccountForUpdate(context.Background(), account.ID)
	require.NoError(t, err)
	checkGotAccount(t, account, account2)
	return account2, err
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1, _ := createRandomAccount(t)
	GetAccount(t, account1)

}

func TestGetAccountForUpdate(t *testing.T) {
	account, _ := createRandomAccount(t)
	GetAccountForUpdate(t, account)

}

func TestDeleteAccount(t *testing.T) {

	account, _ := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account.ID)
	require.Error(t, err)
	require.Empty(t, account2)

}

func TestListAccounts(t *testing.T) {

	// Create 10 random accounts

	var lastAccount Account

	for i := 0; i < 10; i++ {
		lastAccount, _ = createRandomAccount(t)

	}

	// Test listing accounts
	arg := ListAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, lastAccount.Owner, account.Owner)
	}
}

func TestUpdateAccountBalance(t *testing.T) {
	account, _ := createRandomAccount(t)
	arg := UpdateAccountBalanceParams{
		ID:      account.ID,
		Balance: utils.RandomInt63(0, 5000),
	}
	account, err := testQueries.UpdateAccountBalance(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, arg.Balance, account.Balance)
}
func TestAddAccountBalance(t *testing.T) {
	account, _ := createRandomAccount(t)
	arg := AddAccountBalanceParams{
		AccountID: account.ID,
		Amount:    utils.RandomInt63(0, 5000),
	}
	updateAccount, err := testQueries.AddAccountBalance(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, account.Balance+arg.Amount, updateAccount.Balance)
}
