package db

import (
	"context"
	"fmt"
	"simple-bank/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomTestTransfer(fromAccount, toAccount, amount int64) (Transfer, error) {

	args := CreateTransferParams{
		FromAccountID: fromAccount,
		ToAccountID:   toAccount,
		Amount:        amount,
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), args)

	return transfer, err
}

func deleteTestTransfer(id int64) {
	err := testQueries.DeleteAccount(context.Background(), id)

	if err != nil {
		fmt.Println("could not delete transfer ", err)
		return
	}

	fmt.Println("Deleted the transfer record for ", id)
}

func TestCreateTransfer(t *testing.T) {
	fromAccount := CreateRandomTestAccount(t)
	toAccount := CreateRandomTestAccount(t)
	amount := utils.RandomMoney()

	transfer, err := createRandomTestTransfer(fromAccount.ID, toAccount.ID, amount)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.Equal(t, fromAccount.ID, transfer.FromAccountID)
	require.Equal(t, toAccount.ID, transfer.ToAccountID)
	require.Equal(t, amount, transfer.Amount)

	deleteTestTransfer(transfer.ID)
	DeleteTestAccount(fromAccount.ID)
	DeleteTestAccount(toAccount.ID)
}

func TestGetTransfer(t *testing.T) {
	fromAccount := CreateRandomTestAccount(t)
	toAccount := CreateRandomTestAccount(t)
	amount := utils.RandomMoney()

	createTransfer, err := createRandomTestTransfer(fromAccount.ID, toAccount.ID, amount)
	require.NoError(t, err)

	transfer, err := testQueries.GetTransfer(context.Background(), createTransfer.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.Equal(t, fromAccount.ID, transfer.FromAccountID)
	require.Equal(t, toAccount.ID, transfer.ToAccountID)
	require.Equal(t, amount, transfer.Amount)

	deleteTestTransfer(transfer.ID)
	DeleteTestAccount(fromAccount.ID)
	DeleteTestAccount(toAccount.ID)
}

func TestGetTransfers(t *testing.T) {
	fromAccount := CreateRandomTestAccount(t)
	toAccount := CreateRandomTestAccount(t)

	var transferList []Transfer

	for i := 0; i < 10; i++ {
		amount := utils.RandomMoney()

		transfer, _ := createRandomTestTransfer(fromAccount.ID, toAccount.ID, amount)

		transferList = append(transferList, transfer)
	}

	args := GetTransfersParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQueries.GetTransfers(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.Equal(t, fromAccount.ID, transfer.FromAccountID)
		require.Equal(t, toAccount.ID, transfer.ToAccountID)
	}

	for _, transfer := range transferList {
		deleteTestTransfer(transfer.ID)
	}

	DeleteTestAccount(fromAccount.ID)
	DeleteTestAccount(toAccount.ID)
}
