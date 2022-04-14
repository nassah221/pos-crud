package db

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetProductDetails(t *testing.T) {
	arg := ListProductsDetailsByCategoryIDParams{
		CategoryID: 1,
		Limit:      10,
		Offset:     0,
	}
	prods, err := testQueries.ListProductsDetailsByCategoryID(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, prods)
	res := map[string][]interface{}{}
	for _, p := range prods {
		res["products"] = append(res["products"], p)
		if !p.DiscountID.Valid {
			log.Printf("product %v does not have a discount", p.ID)
			log.Printf("%#+v", p)
		}
	}
}
