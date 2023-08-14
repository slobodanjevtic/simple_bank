package db

import (
	"fmt"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	fromAccount := CreateRandomAccount(t)
	toAccount := CreateRandomAccount(t)

	n := 10
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams {
				FromAccountID: fromAccount.ID,
				ToAccountID: toAccount.ID,
				Amount: amount,
			})

			errs <- err
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, fromAccount.ID, transfer.FromAccountID)
		require.Equal(t, toAccount.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromAccount.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntrie(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, toAccount.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntrie(context.Background(), toEntry.ID)
		require.NoError(t, err)

		fromAccountTest := result.FromAccount
		require.NotEmpty(t, fromAccountTest)
		require.Equal(t, fromAccount.ID, fromAccountTest.ID)

		_, err = store.GetAccount(context.Background(), fromAccount.ID)
		require.NoError(t, err)

		toAccountTest := result.ToAccount
		require.NotEmpty(t, toAccountTest)
		require.Equal(t, toAccount.ID, toAccountTest.ID)

		_, err = store.GetAccount(context.Background(), toAccount.ID)
		require.NoError(t, err)
		
		fromDiff := fromAccount.Balance - fromAccountTest.Balance
		toDiff := toAccountTest.Balance - toAccount.Balance
		require.Equal(t, fromDiff, toDiff)

		k := int(fromDiff / amount)
		require.Equal(t, k, i + 1)
	}

	updatedAccountFrom, err := store.GetAccount(context.Background(), fromAccount.ID)
	require.Equal(t, fromAccount.Balance - int64(n) * amount, updatedAccountFrom.Balance)
	require.NoError(t, err)

	updatedAccountTo, err := store.GetAccount(context.Background(), toAccount.ID)
	require.Equal(t, toAccount.Balance + int64(n) * amount, updatedAccountTo.Balance)
	require.NoError(t, err)
}




func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	fromAccount := CreateRandomAccount(t)
	toAccount := CreateRandomAccount(t)

	fmt.Println(">> before transfer: ", fromAccount.Balance, toAccount.Balance)

	n := 10
	amount := int64(10)

	errs := make(chan error)

	for i := 0; i < n / 2; i++ {
		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams {
				FromAccountID: fromAccount.ID,
				ToAccountID: toAccount.ID,
				Amount: amount,
			})

			errs <- err
		}()

		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams {
				FromAccountID: toAccount.ID,
				ToAccountID: fromAccount.ID,
				Amount: amount,
			})

			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	updatedAccountFrom, err := store.GetAccount(context.Background(), fromAccount.ID)
	require.Equal(t, fromAccount.Balance, updatedAccountFrom.Balance)
	require.NoError(t, err)

	updatedAccountTo, err := store.GetAccount(context.Background(), toAccount.ID)
	require.Equal(t, toAccount.Balance, updatedAccountTo.Balance)
	require.NoError(t, err)

	fmt.Println(">> after transfer: ", fromAccount.Balance, toAccount.Balance)
}