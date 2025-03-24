package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zjr71163356/simplebank/utils"
)

func CheckCreatedAccount(t *testing.T, account Account, arg CreateAccountParams) {
	require.NotEmpty(t, account)

	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
}

func CheckGotAccount(t *testing.T, account Account, account2 Account) {
	require.NotEmpty(t, account2)

	require.Equal(t, account.Owner, account2.Owner)
	require.Equal(t, account.Currency, account2.Currency)
	require.Equal(t, account.Balance, account2.Balance)
	require.Equal(t, account.ID, account2.ID)
	require.Equal(t, account.CreatedAt, account2.CreatedAt)
}

func CreateAccount(t *testing.T) (Account, error) {
	arg := CreateAccountParams{
		Owner:    utils.RandomString(5),
		Balance:  utils.RandomInt63(0, 5000),
		Currency: utils.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	CheckCreatedAccount(t, account, arg)
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
	CheckCreatedAccount(t, account, arg)
	return account
}

func GetAccount(t *testing.T, account Account) (Account, error) {

	account2, err := testQueries.GetAccount(context.Background(), account.ID)
	require.NoError(t, err)
	CheckGotAccount(t, account, account2)
	return account2, err
}

func GetAccountForUpdate(t *testing.T, account Account) (Account, error) {
	account2, err := testQueries.GetAccountForUpdate(context.Background(), account.ID)
	require.NoError(t, err)
	CheckGotAccount(t, account, account2)
	return account2, err
}

func TestCreateAccount(t *testing.T) {
	CreateAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1, _ := CreateAccount(t)
	GetAccount(t, account1)

}

func TestGetAccountForUpdate(t *testing.T) {
	account, _ := CreateAccount(t)
	GetAccountForUpdate(t, account)

}

func TestDeleteAccount(t *testing.T) {

	account, _ := CreateAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account.ID)
	require.Error(t, err)
	require.Empty(t, account2)

}

func TestListAccounts(t *testing.T) {

	utils.RandomString(5)
	// Create 10 random accounts
	var createdAccountList []Account

	Owner := utils.RandomOwnerName(1, 5)
	for i := 0; i < 10; i++ {
		account := CreateAccountWithFixedOwner(t, Owner)
		createdAccountList = append(createdAccountList, account)
	}

	// Test listing accounts
	arg := ListAccountsParams{
		Owner:  Owner,
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for idx, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, createdAccountList[idx+5], account)
	}
}

func TestUpdateAccountBalance(t *testing.T) {
	account, _ := CreateAccount(t)
	arg := UpdateAccountBalanceParams{
		ID:      account.ID,
		Balance: utils.RandomInt63(0, 5000),
	}
	account, err := testQueries.UpdateAccountBalance(context.Background(), arg)
	require.NoError(t, err)
	require.NotEqual(t, arg.Balance, account.Balance)
}
