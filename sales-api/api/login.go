package api

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sales-api/constants"
	"sales-api/dto"
	"sales-api/errors"

	"github.com/gin-gonic/gin"
)

var MsgInvalidPassword = fmt.Sprintf("Passcode Not Match")

func (s *Server) GetPassword(ctx *gin.Context) {
	var uri PathIDParam

	if err := ctx.ShouldBindUri(&uri); err != nil {
		log.Printf("[ERR] %v", err)
		errHTTP400(ctx, err)
		return
	}

	cashier, err := s.store.GetCashier(context.Background(), uri.ID)
	if err != nil {
		log.Printf("[ERR] %v", err)
		if err == sql.ErrNoRows {
			errHTTP404(ctx, constants.Cashier)
			return
		}

		errHTTP500(ctx)
		return
	}

	ctx.JSON(http.StatusOK, dto.GenericResponse{
		Success: true,
		Message: "Success",
		Data: map[string]interface{}{
			"passcode": cashier.Password,
		},
	})
}

type loginUserRequestBody struct {
	Password string `json:"passcode" binding:"required,numeric,len=6"`
}

func (s *Server) LoginUser(ctx *gin.Context) {
	var uri PathIDParam
	var body loginUserRequestBody

	if err := ctx.ShouldBindUri(&uri); err != nil {
		log.Printf("[ERR] %v", err)
		errHTTP400(ctx, err)
		return
	}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		log.Printf("[ERR] %v", err)
		vErr, msg := errors.FromFieldValidationErrorPOST(err)
		errHTTP400BodyInvalid(ctx, msg, vErr)
		return
	}

	cashier, err := s.store.GetCashier(context.Background(), uri.ID)
	if err != nil {
		log.Printf("[ERR] %v", err)
		if err == sql.ErrNoRows {
			errHTTP404(ctx, constants.Cashier)
			return
		}

		errHTTP500(ctx)
		return
	}

	if cashier.Password != body.Password {
		ctx.JSON(http.StatusUnauthorized, dto.GenericResponse{
			Success: false,
			Message: MsgInvalidPassword,
		})
		return
	}

	jwt, err := s.tokenMaker.CreateToken(cashier.Name, uri.ID, s.config.JWTDuration)
	if err != nil {
		log.Printf("[ERR] %v", err)
		errHTTP500(ctx)
		return
	}

	ctx.JSON(http.StatusOK, dto.GenericResponse{
		Success: true,
		Message: "Success",
		Data: map[string]interface{}{
			"token": jwt,
		},
	})
}

func (s *Server) LogoutUser(ctx *gin.Context) {
	var uri PathIDParam
	var body loginUserRequestBody

	if err := ctx.ShouldBindUri(&uri); err != nil {
		log.Printf("[ERR] %v", err)
		errHTTP400(ctx, err)
		return
	}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		log.Printf("[ERR] %v", err)
		vErr, msg := errors.FromFieldValidationErrorPOST(err)
		errHTTP400BodyInvalid(ctx, msg, vErr)
		return
	}

	cashier, err := s.store.GetCashier(context.Background(), uri.ID)
	if err != nil {
		log.Printf("[ERR] %v", err)
		if err == sql.ErrNoRows {
			errHTTP404(ctx, constants.Cashier)
			return
		}

		errHTTP500(ctx)
		return
	}

	if cashier.Password != body.Password {
		ctx.JSON(http.StatusUnauthorized, dto.GenericResponse{
			Success: false,
			Message: MsgInvalidPassword,
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.GenericResponse{
		Success: true,
		Message: "Success",
	})
}
