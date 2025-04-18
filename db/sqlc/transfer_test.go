package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zjr71163356/simplebank/utils"
)

func CreateRandomTransfer(t *testing.T) (Transfer, error) {
	fromAccount, err := createRandomAccount(t)
	require.NoError(t, err)

	toAccount, err := createRandomAccount(t)
	require.NoError(t, err)

	arg := CreateTransferParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Amount:        utils.RandomInt63(0, 1000),
	}
	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)
	return transfer, err

}

func CreateTransferWithFixedID(t *testing.T, fromAccountId int64, toAccountId int64) (Transfer, error) {

	arg := CreateTransferParams{
		FromAccountID: fromAccountId,
		ToAccountID:   toAccountId,
		Amount:        utils.RandomInt63(0, 1000),
	}
	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)
	return transfer, err

}

func TestCreateTransfer(t *testing.T) {
	CreateRandomTransfer(t)
}
func TestGetTransfer(t *testing.T) {
	transfer, _ := CreateRandomTransfer(t)
	transfer2, err := testQueries.GetTransfer(context.Background(), transfer.ID)

	require.NoError(t, err)

	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer.ID, transfer2.ID)
	require.Equal(t, transfer.Amount, transfer2.Amount)
	require.Equal(t, transfer.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer.CreatedAt, transfer2.CreatedAt)

}

func TestListTransfers(t *testing.T) {
	var transferList []Transfer
	fromAccount, _ := createRandomAccount(t)
	toAccount, _ := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		transfer, _ := CreateTransferWithFixedID(t, fromAccount.ID, toAccount.ID)
		transferList = append(transferList, transfer)

	}
	arg := ListTransfersParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Limit:         5,
		Offset:        5,
	}

	transferList2, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transferList2)

	for idx, transfer := range transferList2 {
		require.Equal(t, transferList[idx+5], transfer)
	}

}
