package api

import (
	"context"
	"log"
	"net/http"
	"sales-api/dto"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *Server) GetSold(ctx *gin.Context) {
	// disable only full group by sql mode
	if err := s.store.SetSQLMode(context.Background()); err != nil {
		errHTTP500(ctx)
		return
	}

	sold, err := s.store.GetSold(context.Background())
	log.Printf("[ERR] %v", err)
	if err != nil {
		errHTTP500(ctx)
		return
	}
	var resp dto.SoldResponse
	for i := 0; i < len(sold); i++ {
		qty, _ := strconv.Atoi(string(sold[i].TotalQty.([]uint8)))
		amt, _ := strconv.Atoi(string(sold[i].TotalAmount.(([]uint8))))
		v := dto.Sold{
			ProductID:   sold[i].ProductID,
			ProductName: sold[i].ProductName,
			TotalQty:    int32(qty),
			TotalAmount: int64(amt),
		}
		resp.OrderProducts = append(resp.OrderProducts, v)
	}

	ctx.JSON(http.StatusOK, genericResponse{
		Success: true,
		Message: "Success",
		Data:    resp,
	})
}
