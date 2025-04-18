package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransTx(t *testing.T) {
	testStore := NewStore(testDB)

	account1, _ := createRandomAccount(t)
	account2, _ := createRandomAccount(t)
	amount := int64(10)
	n := 5
	fmt.Printf("tx: account1 balance %d account2 balance %d\n", account1.Balance, account2.Balance)

	errs := make(chan error)
	results := make(chan TransferTxResult)
	mark := make(map[int]bool)

	for i := 0; i < n; i++ {
		go func() {
			result, err := testStore.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result

		}()
	}

	for i := 0; i < n; i++ {

		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, amount, transfer.Amount)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.NotEmpty(t, transfer.CreatedAt)
		require.NotEmpty(t, transfer.ID)

		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotEmpty(t, fromEntry.CreatedAt)
		require.NotEmpty(t, fromEntry.ID)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotEmpty(t, toEntry.CreatedAt)
		require.NotEmpty(t, toEntry.ID)

		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		fmt.Printf("account1 balance %d account2 balance %d\n", fromAccount.Balance, toAccount.Balance)
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, mark, k)
		mark[k] = true

	}

	updateFromAccount, err := testStore.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updateToAccount, err := testStore.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance-int64(n)*amount, updateFromAccount.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updateToAccount.Balance)

}
func TestTransTxDeadLock(t *testing.T) {
	testStore := NewStore(testDB)

	account1, _ := createRandomAccount(t)
	account2, _ := createRandomAccount(t)
	amount := int64(10)
	n := 10
	fmt.Printf("tx: account1 balance %d account2 balance %d\n", account1.Balance, account2.Balance)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		var FromAccountID, ToAccountID int64

		if i%2 == 0 {
			FromAccountID = account1.ID
			ToAccountID = account2.ID
		} else {

			FromAccountID = account2.ID
			ToAccountID = account1.ID

		}

		go func() {
			_, err := testStore.TransferTx(context.Background(),  TransferTxParams{
				FromAccountID: FromAccountID,
				ToAccountID:   ToAccountID,
				Amount:        amount,
			})

			errs <- err

		}()
	}

	for i := 0; i < n; i++ {

		err := <-errs
		require.NoError(t, err)

	}

	updateFromAccount, err := testStore.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updateToAccount, err := testStore.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance, updateFromAccount.Balance)
	require.Equal(t, account2.Balance, updateToAccount.Balance)

}
