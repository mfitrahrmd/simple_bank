package db

import (
	"context"
	"database/sql"
	"fmt"
)

type TransferTxParams struct {
	FromAccountID int32   `json:"from_account_id"`
	ToAccountID   int32   `json:"to_account_id"`
	Amount        float64 `json:"amount"`
}
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

type Store interface {
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	Querier
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

// store constructor
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		Queries: New(db),
		db:      db,
	}
}

// wrap query into transactional
// transaction will be rollback if callback fn returns error
func (s *SQLStore) execTx(ctx context.Context, fn func(queries *Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := s.Queries.WithTx(tx)

	err = fn(q)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("transaction error : %v, rollback error : %v", err, rollbackErr)
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit error : %v", err)
	}

	return nil
}

// transfer given amount of money from an account into another
func (s *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var transferResult TransferTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		// create transfer record
		createdTransfer, err := q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		transferResult.Transfer = createdTransfer

		// create transfer entry for sender money
		createdFromEntry, err := s.CreateEntry(ctx, CreateEntryParams{arg.FromAccountID, -arg.Amount})
		if err != nil {
			return err
		}

		transferResult.FromEntry = createdFromEntry

		// create transfer entry for receiver money
		createdToEntry, err := s.CreateEntry(ctx, CreateEntryParams{arg.ToAccountID, arg.Amount})
		if err != nil {
			return err
		}

		transferResult.ToEntry = createdToEntry

		// check and execute transaction first for account with smaller ID to prevent database deadlocks
		if arg.FromAccountID < arg.ToAccountID {
			updatedFromAccount, err := q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:      arg.FromAccountID,
				Balance: -arg.Amount,
			})
			if err != nil {
				return err
			}

			transferResult.FromAccount = updatedFromAccount

			updatedToAccount, err := q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:      arg.ToAccountID,
				Balance: arg.Amount,
			})
			if err != nil {
				return err
			}

			transferResult.ToAccount = updatedToAccount
		} else {
			updatedToAccount, err := q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:      arg.ToAccountID,
				Balance: arg.Amount,
			})
			if err != nil {
				return err
			}

			transferResult.ToAccount = updatedToAccount

			updatedFromAccount, err := q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:      arg.FromAccountID,
				Balance: -arg.Amount,
			})
			if err != nil {
				return err
			}

			transferResult.FromAccount = updatedFromAccount
		}

		return nil
	})
	if err != nil {
		return TransferTxResult{}, err
	}

	return transferResult, nil
}
