package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStore_TransferTx(t *testing.T) {
	s := NewStore(dbTest)

	sender := testCreateAccount(t)
	receiver := testCreateAccount(t)

	arg := TransferTxParams{
		FromAccountID: sender.ID,
		ToAccountID:   receiver.ID,
		Amount:        20,
	}

	results := make(chan TransferTxResult)
	errs := make(chan error)

	for i := 0; i < 5; i++ {
		go func() {
			createdTransfer, err := s.TransferTx(context.Background(), arg)
			errs <- err
			results <- createdTransfer
		}()
	}

	for i := 0; i < 5; i++ {
		e := <-errs
		require.NoError(t, e)

		r := <-results
		require.NotEmpty(t, r)

		transfer := r.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, sender.ID, transfer.FromAccountID)
		require.Equal(t, receiver.ID, transfer.ToAccountID)
		require.Equal(t, arg.Amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err := s.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		fromEntry := r.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, sender.ID, fromEntry.AccountID)
		require.Equal(t, -arg.Amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)
		_, err = s.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := r.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, receiver.ID, toEntry.AccountID)
		require.Equal(t, arg.Amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)
		_, err = s.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		fromAccount := r.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, sender.ID, fromAccount.ID)

		toAccount := r.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, receiver.ID, toAccount.ID)

		diff1 := sender.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - receiver.Balance
		require.Equal(t, diff1, diff2)
	}

	senderAccountTx, err := queriesTest.GetAccount(context.Background(), sender.ID)
	require.NoError(t, err)
	require.Equal(t, senderAccountTx.Balance, sender.Balance-(arg.Amount*5))

	receiverAccountTx, err := queriesTest.GetAccount(context.Background(), receiver.ID)
	require.NoError(t, err)
	require.Equal(t, receiverAccountTx.Balance, receiver.Balance+(arg.Amount*5))
}
