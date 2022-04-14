package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateCashier(t *testing.T) {
	arg := CreateCashierParams{
		Name:     "Tom",
		Password: "12345",
	}

	cashier, err := testQueries.CreateCashier(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, cashier)

	insertID, err := cashier.LastInsertId()
	require.NoError(t, err)

	checkCashier, _ := testQueries.GetCashier(context.Background(), int32(insertID))
	require.Equal(t, arg.Name, checkCashier.Name)
	require.Equal(t, arg.Password, checkCashier.Password)
}

func TestGetCashier(t *testing.T) {

}
