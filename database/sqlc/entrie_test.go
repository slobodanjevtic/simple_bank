package db

import (
	"database/sql"
	"context"
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"

	"github.com/stretchr/testify/require"
)

func CreateRandomEntry(t *testing.T) Entry {
	amount, err := faker.RandomInt(0, 1000, 1)
	randAccount := CreateRandomAccount(t)
	arg := CreateEntrieParams {
		AccountID: randAccount.ID,
		Amount: int64(amount[0]),
	}
	
	entry, err := testQueries.CreateEntrie(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry;
}

func TestCreateEntry(t *testing.T) {
	CreateRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	randEntry := CreateRandomEntry(t)
	entry, err := testQueries.GetEntrie(context.Background(), randEntry.ID)
	
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, randEntry.ID, entry.ID)
	require.Equal(t, randEntry.AccountID, entry.AccountID)
	require.Equal(t, randEntry.Amount, entry.Amount)
	require.WithinDuration(t, randEntry.CreatedAt, entry.CreatedAt, time.Second)
}

func TestUpdateEntry(t *testing.T) {
	randEntry := CreateRandomEntry(t)
	amount, err := faker.RandomInt(0, 1000, 1)

	arg := UpdateEntrieParams {
		ID: randEntry.ID,
		Amount: int64(amount[0]),
	}
	
	entry, err := testQueries.UpdateEntrie(context.Background(), arg)
	
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, randEntry.ID, entry.ID)
	require.Equal(t, randEntry.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.WithinDuration(t, randEntry.CreatedAt, entry.CreatedAt, time.Second)
}

func TestDeleteEntry(t *testing.T) {
	randEntry := CreateRandomEntry(t)
	err := testQueries.DeleteEntrie(context.Background(), randEntry.ID)
	
	require.NoError(t, err)

	entry, err := testQueries.GetEntrie(context.Background(), randEntry.ID)
	
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, entry)
}


func TestListEntrys(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomEntry(t)
	}
	
	arg := ListEntriesParams {
		Limit: 5,
		Offset: 5,
	}

	Entrys, err := testQueries.ListEntries(context.Background(), arg)
	
	require.NoError(t, err)
	require.Len(t, Entrys, 5)
	for _, Entry := range Entrys {
		require.NotEmpty(t, Entry)
	}
}