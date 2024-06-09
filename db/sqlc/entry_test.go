package db

import (
	"context"
	"fmt"
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

func DeleteRandomTestEntry(id int64) {
	err := testQueries.DeleteAccount(context.Background(), id)

	if err != nil {
		fmt.Println("could not delete entry ", err)
		return
	}

	fmt.Println("Deleted the entry record for ", id)
}

func TestCreateEntry(t *testing.T) {
	account := CreateRandomTestAccount(t)
	entry := CreateRandomTestEntry(t, account.ID)
	DeleteRandomTestEntry(entry.ID)
	DeleteTestAccount(account.ID)
}

func TestGetEntry(t *testing.T) {
	account := CreateRandomTestAccount(t)
	entry1 := CreateRandomTestEntry(t, account.ID)

	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)

	DeleteRandomTestEntry(entry1.ID)
	DeleteTestAccount(account.ID)
}

func TestGetEntries(t *testing.T) {
	account := CreateRandomTestAccount(t)
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

	for _, entries := range entriesList {
		DeleteRandomTestEntry(entries.ID)
	}

	DeleteTestAccount(account.ID)
}
