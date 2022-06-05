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

func (s *Server) ListCategories(ctx *gin.Context) { //nolint
	var query PaginationQueryParams

	if err := ctx.ShouldBindQuery(&query); err != nil {
		log.Printf("[ERR] %v", err)
		errHTTP400(ctx, err)
		return
	}

	arg := db.ListCategoriesParams{
		Limit:  query.Limit,
		Offset: query.Skip,
	}
	c, err := s.store.ListCategories(context.Background(), arg)
	if err != nil {
		errHTTP500(ctx)
		return
	}

	resp := dto.ListCategory{
		Meta: dto.Meta{
			Limit: int(query.Limit),
			Skip:  int(query.Skip),
			Total: len(c),
		},
	}

	var categories []dto.Category

	for i := 0; i < len(c); i++ {
		categories = append(categories, dto.Category{
			ID:   c[i].ID,
			Name: c[i].Name,
		})
	}
	resp.Category = categories

	ctx.JSON(http.StatusOK, dto.GenericResponse{
		Success: true,
		Message: "Success",
		Data:    resp,
	})
}

func (s *Server) GetCategory(ctx *gin.Context) {
	var uri PathIDParam

	if err := ctx.ShouldBindUri(&uri); err != nil {
		log.Printf("[ERR] %v", err)
		errHTTP400(ctx, err)
		return
	}

	c, err := s.store.GetCategory(context.Background(), uri.ID)
	if err != nil {
		log.Printf("[ERR] %v", err)
		if err == sql.ErrNoRows {
			errHTTP404(ctx, constants.Category)
			return
		}

		errHTTP500(ctx)
		return
	}

	ctx.JSON(http.StatusOK, dto.GenericResponse{
		Success: true,
		Message: "Success",
		Data: dto.Category{
			ID:   c.ID,
			Name: c.Name,
		},
	})
}

// Is also used for PUT because of the body being the same
type createCategoryRequestBody struct {
	Name string `json:"name" binding:"required"`
}

func (s *Server) CreateCategory(ctx *gin.Context) {
	var body createCategoryRequestBody

	if err := ctx.ShouldBindJSON(&body); err != nil {
		log.Printf("[ERR] %v", err)
		vErr, msg := errors.FromFieldValidationErrorPOST(err)
		errHTTP400BodyInvalid(ctx, msg, vErr)
		return
	}

	c, err := s.store.CreateCategoryWithReturn(context.Background(), body.Name)
	if err != nil {
		errHTTP500(ctx)
		return
	}

	ctx.JSON(http.StatusOK, dto.GenericResponse{
		Success: true,
		Message: "Success",
		Data: dto.Category{
			ID:        c.ID,
			Name:      c.Name,
			UpdatedAt: &c.UpdatedAt,
			CreatedAt: &c.CreatedAt,
		},
	})
}

func (s *Server) UpdateCategory(ctx *gin.Context) {
	var uri PathIDParam
	var body createCategoryRequestBody

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

	_, err := s.store.GetCategory(context.Background(), uri.ID)
	if err != nil {
		log.Printf("[ERR] %v", err)
		if err == sql.ErrNoRows {
			errHTTP404(ctx, constants.Category)
			return
		}

		errHTTP500(ctx)
		return
	}

	arg := db.UpdateCategoryParams{
		Name: body.Name,
		ID:   uri.ID,
	}
	if err := s.store.UpdateCategory(context.Background(), arg); err != nil {
		errHTTP500(ctx)
		return
	}

	ctx.JSON(http.StatusOK, dto.GenericResponse{
		Success: true,
		Message: "Success",
	})
}

func (s *Server) DeleteCategory(ctx *gin.Context) { //nolint
	var uri PathIDParam
	if err := ctx.ShouldBindUri(&uri); err != nil {
		log.Printf("[ERR] %v", err)
		errHTTP400(ctx, err)
		return
	}

	_, err := s.store.GetCategory(context.Background(), uri.ID)
	if err != nil {
		log.Printf("[ERR] %v", err)
		if err == sql.ErrNoRows {
			errHTTP404(ctx, constants.Category)
			return
		}

		errHTTP500(ctx)
		return
	}

	if err := s.store.DeleteCategory(context.Background(), uri.ID); err != nil {
		errHTTP500(ctx)
		return
	}

	ctx.JSON(http.StatusOK, dto.GenericResponse{
		Success: true,
		Message: "Success",
	})
}
