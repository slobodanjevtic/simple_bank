package db

import (
	"database/sql"
	"context"
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"

	"github.com/stretchr/testify/require"
)

func CreateRandomTransfer(t *testing.T) Transfer {
	amount, err := faker.RandomInt(0, 1000, 1)
	randAccountFrom := CreateRandomAccount(t)
	randAccountTo := CreateRandomAccount(t)

	arg := CreateTransferParams {
		FromAccountID: randAccountFrom.ID,
		ToAccountID: randAccountTo.ID,
		Amount: int64(amount[0]),
	}
	
	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	
	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer;
}

func TestCreateTransfer(t *testing.T) {
	CreateRandomTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	randTransfer := CreateRandomTransfer(t)
	transfer, err := testQueries.GetTransfer(context.Background(), randTransfer.ID)
	
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, randTransfer.ID, transfer.ID)
	require.Equal(t, randTransfer.FromAccountID, transfer.FromAccountID)
	require.Equal(t, randTransfer.ToAccountID, transfer.ToAccountID)
	require.Equal(t, randTransfer.Amount, transfer.Amount)
	require.WithinDuration(t, randTransfer.CreatedAt, transfer.CreatedAt, time.Second)
}

func TestUpdateTransfer(t *testing.T) {
	randTransfer := CreateRandomTransfer(t)
	amount, err := faker.RandomInt(0, 1000, 1)

	arg := UpdateTransferParams {
		ID: randTransfer.ID,
		Amount: int64(amount[0]),
	}
	
	transfer, err := testQueries.UpdateTransfer(context.Background(), arg)
	
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, randTransfer.ID, transfer.ID)
	require.Equal(t, randTransfer.FromAccountID, transfer.FromAccountID)
	require.Equal(t, randTransfer.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)
	require.WithinDuration(t, randTransfer.CreatedAt, transfer.CreatedAt, time.Second)
}

func TestDeleteTransfer(t *testing.T) {
	randTransfer := CreateRandomTransfer(t)
	err := testQueries.DeleteTransfer(context.Background(), randTransfer.ID)
	
	require.NoError(t, err)

	transfer, err := testQueries.GetTransfer(context.Background(), randTransfer.ID)
	
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, transfer)
}


func TestListTransfers(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomTransfer(t)
	}
	
	arg := ListTransfersParams {
		Limit: 5,
		Offset: 5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	
	require.NoError(t, err)
	require.Len(t, transfers, 5)
	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}