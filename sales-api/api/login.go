package api

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sales-api/dto"

	"github.com/gin-gonic/gin"
)

var MsgInvalidPassword = fmt.Sprintf("Incorrect Password")

type loginUserRequestURI struct {
	ID int `uri:"id" binding:"required,min=1"`
}
type loginUserRequestBody struct {
	Password string `json:"passcode" binding:"required,numeric,len=6"`
}

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
			errHTTP404(ctx)
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
		errHTTP400(ctx, err)
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

	if cashier.Password != body.Password {
		ctx.JSON(http.StatusUnauthorized, dto.GenericResponse{
			Success: false,
			Message: MsgInvalidPassword,
		})
		return
	}

	jwt, err := s.tokenMaker.CreateToken(cashier.Name, uri.ID, tokenDuration)
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
