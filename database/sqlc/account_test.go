package db

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"testing"
)

func testCreateAccount(t *testing.T) Account {
	mockCreateAccount := CreateAccountParams{
		Owner:    "Rama",
		Balance:  100,
		Currency: "USD",
	}

	createdAccount, err := queriesTest.CreateAccount(context.Background(), mockCreateAccount)
	require.NoError(t, err)
	require.NotEmpty(t, createdAccount)
	require.Equal(t, createdAccount.Owner, mockCreateAccount.Owner)
	require.Equal(t, createdAccount.Balance, mockCreateAccount.Balance)
	require.Equal(t, createdAccount.Currency, mockCreateAccount.Currency)
	require.NotZero(t, createdAccount.ID)
	require.NotZero(t, createdAccount.CreatedAt)

	return createdAccount
}

func TestQueries_CreateAccount(t *testing.T) {
	testCreateAccount(t)
}

func TestQueries_GetAccount(t *testing.T) {
	createdAccount := testCreateAccount(t)

	foundAccount, err := queriesTest.GetAccount(context.Background(), createdAccount.ID)
	require.NoError(t, err)
	require.NotEmpty(t, foundAccount)
	require.Equal(t, createdAccount.ID, foundAccount.ID)
	require.Equal(t, createdAccount.Owner, foundAccount.Owner)
	require.Equal(t, createdAccount.Balance, foundAccount.Balance)
	require.Equal(t, createdAccount.Currency, foundAccount.Currency)
}

func TestQueries_UpdateAccount(t *testing.T) {
	createdAccount := testCreateAccount(t)

	arg := UpdateAccountParams{
		ID:      createdAccount.ID,
		Balance: 1200,
	}

	updatedAccount, err := queriesTest.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, createdAccount.ID, updatedAccount.ID)
	require.Equal(t, createdAccount.Owner, updatedAccount.Owner)
	require.Equal(t, arg.Balance, updatedAccount.Balance)
	require.Equal(t, createdAccount.Currency, updatedAccount.Currency)
}

func TestQueries_DeleteAccount(t *testing.T) {
	createdAccount := testCreateAccount(t)

	err := queriesTest.DeleteAccount(context.Background(), createdAccount.ID)
	require.NoError(t, err)

	foundAccount, err := queriesTest.GetAccount(context.Background(), createdAccount.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, foundAccount)
}

func TestQueries_ListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		testCreateAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	foundAccounts, err := queriesTest.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, foundAccounts, 5)

	for _, account := range foundAccounts {
		require.NotEmpty(t, account)
	}
}
