package db

import (
	"context"
	"database/sql"
	"fmt"
	"simple-bank/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func CreateRandomTestAccount(t *testing.T) Account {
	args := CreateAccountParams{
		Owner:    utils.RandomOwner(),
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, args.Owner, account.Owner)
	require.Equal(t, args.Balance, account.Balance)
	require.Equal(t, args.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func DeleteTestAccount(id int64) {
	err := testQueries.DeleteAccount(context.Background(), id)

	if err != nil {
		fmt.Println("could not delete account ", err)
		return
	}

	fmt.Println("Deleted the account record for ", id)
}

func TestCreateAccount(t *testing.T) {
	account := CreateRandomTestAccount(t)
	DeleteTestAccount(account.ID)
}

func TestGetAccount(t *testing.T) {
	account1 := CreateRandomTestAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)

	DeleteTestAccount(account1.ID)

}

func TestUpdateAccount(t *testing.T) {
	account1 := CreateRandomTestAccount(t)

	args := UpdateAccountParams{
		ID:      account1.ID,
		Balance: utils.RandomMoney(),
	}

	account2, err := testQueries.UpdateAccount(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, args.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)

	DeleteTestAccount(account1.ID)

}

func TestDeleteAccount(t *testing.T) {
	account1 := CreateRandomTestAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestListAccount(t *testing.T) {
	var accList []Account
	for i := 1; i <= 10; i++ {
		accList = append(accList, CreateRandomTestAccount(t))
	}

	args := ListAccountParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccount(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)
	require.Len(t, accounts, 5)

	for _, acc := range accounts {
		require.NotEmpty(t, acc)
	}

	for _, acc := range accList {
		DeleteTestAccount(acc.ID)
	}
}
