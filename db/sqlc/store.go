package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all the function to execute db queries and transactions
type Store interface {
	Querier
	TransferTxn(ctx context.Context, args TransferTxnParam) (TransferTxnResult, error)
}

// SQLStore provides all the function to execute SQL queries and transactions
type SQLStore struct {
	db *sql.DB
	*Queries
}

// creates a new store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

func (store *SQLStore) execTxn(ctx context.Context, fn func(*Queries) error) error {
	txn, err := store.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	q := New(txn)
	err = fn(q)

	if err != nil {
		if rbErr := txn.Rollback(); rbErr != nil {
			return fmt.Errorf("txn Err: %v, rollBack Err: %v", err, rbErr)
		}
		return err
	}

	return txn.Commit()
}

type TransferTxnParam struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxnResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func getTransferParam(args TransferTxnParam) *CreateTransferParams {
	return &CreateTransferParams{
		FromAccountID: args.FromAccountID,
		ToAccountID:   args.ToAccountID,
		Amount:        args.Amount,
	}
}

// TransferTxn performs the transfer of amount between two accounts
// It creates the transfer record, add account entries, and update the accounts' balance in a single transaction
func (store *SQLStore) TransferTxn(ctx context.Context, args TransferTxnParam) (TransferTxnResult, error) {
	var result TransferTxnResult

	err := store.execTxn(ctx, func(q *Queries) error {
		var err error
		result.Transfer, err = q.CreateTransfer(ctx, *getTransferParam(args))

		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.FromAccountID,
			Amount:    -args.Amount,
		})

		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.ToAccountID,
			Amount:    args.Amount,
		})

		if err != nil {
			return err
		}

		if args.FromAccountID < args.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, args.FromAccountID, -args.Amount, args.ToAccountID, args.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, args.ToAccountID, args.Amount, args.FromAccountID, -args.Amount)
		}
		if err != nil {
			return err
		}
		return nil
	})

	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (fromAccount Account, toAccount Account, err error) {
	fromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	toAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	if err != nil {
		return
	}

	return
}
