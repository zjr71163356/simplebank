package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zjr71163356/simplebank/utils"
)

func CreateRandomEntry(t *testing.T) (Entry, error) {
	account, _ := CreateRandomAccount(t)

	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    utils.RandomInt63(-1000, 1000),
	}
	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry, err

}

func CreateRandomEntryWithFixedAccount(t *testing.T, account Account) (Entry, error) {
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    utils.RandomInt63(-1000, 1000),
	}
	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry, err
}

func GetRandomEntry(t *testing.T) {
	entry1, _ := CreateRandomEntry(t)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)

}

func TestCreateEntry(t *testing.T) {
	CreateRandomEntry(t)
}

func TestEntry(t *testing.T) {
	GetRandomEntry(t)
}

func TestListEntries(t *testing.T) {
	var entrylist []Entry
	account, _ := CreateRandomAccount(t)

	for i := 0; i < 10; i++ {
		entry, _ := CreateRandomEntryWithFixedAccount(t, account)
		entrylist = append(entrylist, entry)
	}

	arg := ListEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for idx, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, entrylist[idx+5], entry)
	}

}
