package db

import (
	"context"
	"simple-bank/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func CreateRandomTestEntry(t *testing.T, accountId int64) Entry {
	amount := utils.RandomMoney()

	args := CreateEntryParams{
		AccountID: accountId,
		Amount:    amount,
	}

	entry, err := testQueries.CreateEntry(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, accountId, entry.AccountID)
	require.Equal(t, amount, entry.Amount)

	return entry
}

func TestCreateEntry(t *testing.T) {
	account := createRandomTestAccount(t)
	CreateRandomTestEntry(t, account.ID)
}

func TestGetEntry(t *testing.T) {
	account := createRandomTestAccount(t)
	entry1 := CreateRandomTestEntry(t, account.ID)

	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

func TestGetEntries(t *testing.T) {
	account := createRandomTestAccount(t)
	var entriesList []Entry
	for i := 0; i < 10; i++ {
		entriesList = append(entriesList, CreateRandomTestEntry(t, account.ID))
	}

	args := GetEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	}

	entries, err := testQueries.GetEntries(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}
