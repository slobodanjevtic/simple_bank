package db

import (
	"database/sql"
	"context"
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/require"
)

func CreateRandomAccount(t *testing.T) Account {
	balance, err := faker.RandomInt(0, 1000, 1)
	arg := CreateAccountParams {
		Owner: faker.FirstName(),
		Balance: int64(balance[0]),
		Currency: faker.Currency(),
	}
	
	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account;
}

func TestCreateAccount(t *testing.T) {
	CreateRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	randAccount := CreateRandomAccount(t)
	account, err := testQueries.GetAccount(context.Background(), randAccount.ID)
	
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, randAccount.ID, account.ID)
	require.Equal(t, randAccount.Owner, account.Owner)
	require.Equal(t, randAccount.Balance, account.Balance)
	require.Equal(t, randAccount.Currency, account.Currency)
	require.WithinDuration(t, randAccount.CreatedAt, account.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	randAccount := CreateRandomAccount(t)
	balance, err := faker.RandomInt(0, 1000, 1)

	arg := UpdateAccountParams {
		ID: randAccount.ID,
		Balance: int64(balance[0]),
	}
	
	account, err := testQueries.UpdateAccount(context.Background(), arg)
	
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, randAccount.ID, account.ID)
	require.Equal(t, randAccount.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, randAccount.Currency, account.Currency)
	require.WithinDuration(t, randAccount.CreatedAt, account.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	randAccount := CreateRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), randAccount.ID)
	
	require.NoError(t, err)

	account, err := testQueries.GetAccount(context.Background(), randAccount.ID)
	
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account)
}


func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomAccount(t)
	}
	
	arg := ListAccountsParams {
		Limit: 5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	
	require.NoError(t, err)
	require.Len(t, accounts, 5)
	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}