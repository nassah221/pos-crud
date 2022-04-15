package api

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	db "sales-api/db/sqlc"
	"sales-api/dto"
	"sales-api/errors"

	"github.com/gin-gonic/gin"
)

var (
	MsgNotFound       = fmt.Sprintf("Resource Not Found")
	MsgInvalidRequest = fmt.Sprintf("Unable to Validate Request")
)

var (
	ErrInternalServer = fmt.Sprintln("Something went wrong")
)

type createCashierRequest struct {
	Name     string `json:"name" binding:"required"`
	Passcode string `json:"passcode" binding:"required,numeric,len=6"`
}

func (s *Server) CreateCashier(ctx *gin.Context) {
	var req createCashierRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[ERR] %v", err)

		vErr, msg := errors.FromFieldValidationErrorPOST(err)
		ctx.JSON(http.StatusBadRequest, dto.GenericResponse{
			Success: false,
			Message: msg,
			Error:   vErr,
		})
		return
	}

	arg := db.CreateCashierParams{Name: req.Name, Password: req.Passcode}

	cashier, err := s.store.CreateCashierWithReturn(context.Background(), arg)
	if err != nil {
		log.Printf("[ERR] %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.GenericResponse{
			Success: false,
			Error:   ErrInternalServer,
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.GenericResponse{
		Success: true,
		Message: "Success",
		Data: dto.Cashier{
			ID:        int(cashier.ID),
			Name:      cashier.Name,
			Passcode:  cashier.Password,
			UpdatedAt: &cashier.UpdatedAt,
			CreatedAt: &cashier.CreatedAt,
		},
	})
}

func (s *Server) GetCashier(ctx *gin.Context) {
	var uri PathIDParam

	if err := ctx.ShouldBindUri(&uri); err != nil {
		log.Printf("[ERR] %v", err)
		ctx.JSON(http.StatusBadRequest, dto.GenericResponse{
			Success: false,
			Message: MsgInvalidRequest,
			Error:   err.Error(),
		})
		return
	}

	cashier, err := s.store.GetCashier(context.Background(), uri.ID)
	if err != nil {
		log.Printf("[ERR] %v", err)
		if err == sql.ErrNoRows {
			errHTTP404(ctx)
			return
		}

		errHTTP500(ctx)
		return
	}

	ctx.JSON(http.StatusOK, dto.GenericResponse{
		Success: true,
		Message: "Success",
		Data: dto.Cashier{
			ID:   int(uri.ID),
			Name: cashier.Name,
		},
	})
}

func (s *Server) ListCashiers(ctx *gin.Context) { //nolint
	var query PaginationQueryParams

	if err := ctx.ShouldBindQuery(&query); err != nil {
		log.Printf("[ERR] %v", err)
		ctx.JSON(http.StatusBadRequest, dto.GenericResponse{
			Success: false,
			Message: MsgInvalidRequest,
			Error:   err.Error(),
		})
		return
	}

	arg := db.ListCashiersParams{
		Limit:  query.Limit,
		Offset: query.Skip,
	}

	cashiers, err := s.store.ListCashiers(context.Background(), arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.GenericResponse{
			Success: false,
			Error:   ErrInternalServer,
		})
		return
	}

	var cashiersList []dto.Cashier
	resp := dto.ListCashiers{
		Meta: dto.Meta{
			Limit: int(query.Limit),
			Skip:  int(query.Skip),
			Total: len(cashiers),
		},
	}

	for i := 0; i < len(cashiers); i++ {
		cashiersList = append(cashiersList, dto.Cashier{
			ID:   int(cashiers[i].ID),
			Name: cashiers[i].Name,
		})
	}
	resp.Cashiers = cashiersList

	ctx.JSON(http.StatusOK, dto.GenericResponse{
		Success: true,
		Message: "Success",
		Data:    resp,
	})
}

type updateCashierRequest struct {
	Name     string `json:"name,omitempty" binding:"required"`
	Passcode string `json:"passcode,omitempty" binding:"numeric,len=6"`
}

func (s *Server) UpdateCashier(ctx *gin.Context) {
	var req updateCashierRequest
	var uri PathIDParam

	if err := ctx.ShouldBindUri(&uri); err != nil {
		log.Printf("[ERR] %v", err)
		errHTTP400(ctx, err)
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[ERR] %v", err)
		vErr, msg := errors.FromFieldValidationErrorPUT(err)
		errHTTP400BodyInvalid(ctx, msg, vErr)
		return
	}

	cashier, err := s.store.GetCashier(context.Background(), uri.ID)
	if err != nil {
		log.Printf("[ERR] %v", err)
		if err == sql.ErrNoRows {
			errHTTP404(ctx)
			return
		}
		errHTTP500(ctx)
		return
	}

	if req.Name == "" {
		req.Name = cashier.Name
	}
	if req.Passcode == "" {
		req.Passcode = cashier.Password
	}

	arg := db.UpdateCashierParams{Name: req.Name, Password: req.Passcode, ID: uri.ID}
	if err := s.store.UpdateCashier(context.Background(), arg); err != nil {
		log.Printf("[ERR] %v", err)
		errHTTP500(ctx)
		return
	}

	ctx.JSON(http.StatusOK, dto.GenericResponse{
		Success: true,
		Message: "Success",
	})
}

func (s *Server) DeleteCashier(ctx *gin.Context) {
	var uri PathIDParam

	if err := ctx.ShouldBindUri(&uri); err != nil {
		log.Printf("[ERR] %v", err)
		errHTTP400(ctx, err)
		return
	}

	_, err := s.store.GetCashier(context.Background(), uri.ID)
	if err != nil {
		log.Printf("[ERR] %v", err)
		if err == sql.ErrNoRows {
			errHTTP404(ctx)
			return
		}

		errHTTP500(ctx)
		return
	}

	if err := s.store.DeleteCashier(context.Background(), uri.ID); err != nil {
		log.Printf("[ERR] %v", err)
		errHTTP500(ctx)
		return
	}

	ctx.JSON(http.StatusOK, dto.GenericResponse{
		Success: true,
		Message: "Success",
	})
}
