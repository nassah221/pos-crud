package api

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"sales-api/constants"
	db "sales-api/db/sqlc"
	"sales-api/dto"
	"sales-api/errors"

	"github.com/gin-gonic/gin"
)

// todo: figure out what &subtotal query param in the docs is and implement it

type listPaymentQueryParams struct {
	Subtotal int `form:"subtotal"`
	PaginationQueryParams
}

func (s *Server) ListPayments(ctx *gin.Context) {
	var query listPaymentQueryParams

	if err := ctx.ShouldBindQuery(&query); err != nil {
		log.Printf("[ERR] %v", err)
		errHTTP400(ctx, err)
		return
	}

	arg := db.ListPaymentsParams{
		Limit:  query.Limit,
		Offset: query.Skip,
	}

	payments, err := s.store.ListPayments(context.Background(), arg)
	if err != nil {
		errHTTP500(ctx)
		return
	}

	resp := dto.ListPayment{
		Meta: dto.Meta{
			Total: len(payments),
			Limit: int(query.Limit),
			Skip:  int(query.Skip),
		},
	}
	var p []dto.Payment
	for i := 0; i < len(payments); i++ {
		p = append(p, dto.Payment{
			ID:   int(payments[i].ID),
			Name: payments[i].Name,
		})
	}
	resp.Payment = p

	ctx.JSON(http.StatusOK, dto.GenericResponse{
		Success: true,
		Message: "Success",
		Data:    resp,
	})
}
func (s *Server) GetPayment(ctx *gin.Context) {
	var uri PathIDParam

	if err := ctx.ShouldBindUri(&uri); err != nil {
		log.Printf("[ERR] %v", err)
		errHTTP400(ctx, err)
		return
	}

	payment, err := s.store.GetPayment(context.Background(), uri.ID)
	if err != nil {
		log.Printf("[ERR] %v", err)
		if err == sql.ErrNoRows {
			errHTTP404(ctx, constants.Payment)
			return
		}

		errHTTP500(ctx)
		return
	}

	ctx.JSON(http.StatusOK, dto.GenericResponse{
		Success: true,
		Message: "Success",
		Data: dto.Payment{
			ID:   int(payment.ID),
			Name: payment.Name,
			Type: payment.Type,
			Logo: payment.Logo.String,
		},
	})
}

type createPaymentRequestBody struct {
	Name string `json:"name" binding:"required"`
	Type string `json:"type" binding:"required,oneof=CASH E-WALLET EDC"`
	Logo string `json:"logo"`
}

func (s *Server) CreatePayment(ctx *gin.Context) {
	var body createPaymentRequestBody

	if err := ctx.ShouldBindJSON(&body); err != nil {
		log.Printf("[ERR] %v", err)
		vErr, msg := errors.FromFieldValidationErrorPOST(err)
		errHTTP400BodyInvalid(ctx, msg, vErr)
		return
	}

	arg := db.CreatePaymentParams{Name: body.Name, Type: body.Type, Logo: sql.NullString{body.Logo, body.Logo != ""}}
	payment, err := s.store.CreatePaymentWithReturn(context.Background(), arg)
	if err != nil {
		log.Printf("[ERR] %v", err)
		errHTTP500(ctx)
		return
	}

	ctx.JSON(http.StatusOK, dto.GenericResponse{
		Success: true,
		Message: "Success",
		Data: dto.Payment{
			ID:        int(payment.ID),
			Name:      payment.Name,
			Type:      payment.Type,
			Logo:      payment.Logo.String,
			CreatedAt: &payment.CreatedAt,
			UpdatedAt: &payment.UpdatedAt,
		},
	})
}

type updatePaymentRequestBody struct {
	Name string `json:"name"`
	Type string `json:"type" binding:"oneof=CASH E-WALLET EDC"`
	Logo string `json:"logo"`
}

func (s *Server) UpdatePayment(ctx *gin.Context) {
	var uri PathIDParam
	var body updatePaymentRequestBody

	if err := ctx.ShouldBindUri(&uri); err != nil {
		log.Printf("[ERR] %v", err)
		errHTTP400(ctx, err)
		return
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		log.Printf("[ERR] %v", err)
		vErr, msg := errors.FromFieldValidationErrorPUT(err)
		errHTTP400BodyInvalid(ctx, msg, vErr)
		return
	}

	payment, err := s.store.GetPayment(context.Background(), uri.ID)
	if err != nil {
		log.Printf("[ERR] %v", err)
		if err == sql.ErrNoRows {
			errHTTP404(ctx, constants.Payment)
			return
		}

		errHTTP500(ctx)
		return
	}

	arg := db.UpdatePaymentParams{
		Name: body.Name,
		Type: body.Type,
		Logo: sql.NullString{body.Logo, body.Logo != ""},
		ID:   uri.ID,
	}

	if body.Name == "" {
		arg.Name = body.Name
	}
	if body.Type == "" {
		arg.Type = payment.Type
	}
	if body.Logo == "" {
		arg.Logo = sql.NullString{body.Logo, body.Logo != ""}
	}

	if err := s.store.UpdatePayment(context.Background(), arg); err != nil {
		log.Printf("[ERR] %v", err)
		errHTTP500(ctx)
		return
	}

	ctx.JSON(http.StatusOK, dto.GenericResponse{
		Success: true,
		Message: "Success",
	})
}
func (s *Server) DeletePayment(ctx *gin.Context) {
	var uri PathIDParam
	if err := ctx.ShouldBindUri(&uri); err != nil {
		log.Printf("[ERR] %v", err)
		errHTTP400(ctx, err)
		return
	}

	_, err := s.store.GetPayment(context.Background(), uri.ID)
	if err != nil {
		log.Printf("[ERR] %v", err)
		if err == sql.ErrNoRows {
			errHTTP404(ctx, constants.Payment)
			return
		}

		errHTTP500(ctx)
		return
	}

	if err := s.store.DeletePayment(context.Background(), uri.ID); err != nil {
		log.Printf("[ERR] %v", err)
		errHTTP500(ctx)
		return
	}

	ctx.JSON(http.StatusOK, dto.GenericResponse{
		Success: true,
		Message: "Success",
	})
}
