package api

import (
	"context"
	"log"
	"net/http"
	"sales-api/dto"

	"github.com/gin-gonic/gin"
)

func (s *Server) GetRevenue(ctx *gin.Context) {
	r, err := s.store.GetRevenue(context.Background())
	if err != nil {
		log.Printf("[ERR] %v", err)
		errHTTP500(ctx)
		return
	}

	var totalRevenue float64
	payments := []dto.RevenuePayment{}
	for i := 0; i < len(r); i++ {
		payments = append(payments, dto.RevenuePayment{
			ID:          r[i].PaymentID,
			Name:        r[i].Name,
			Logo:        r[i].Logo.String,
			TotalAmount: float64(r[i].TotalAmount),
		})
		totalRevenue += float64(r[i].TotalAmount)
	}

	ctx.JSON(http.StatusOK, genericResponse{
		Success: true,
		Message: "Success",
		Data: dto.Revenue{
			TotalRevenue: totalRevenue,
			PaymentType:  payments,
		},
	})
}
