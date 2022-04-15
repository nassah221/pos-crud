package api

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sales-api/constants"
	db "sales-api/db/sqlc"
	"sales-api/dto"
	"sales-api/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *Server) CreateProduct(ctx *gin.Context) {
	var req dto.ProductCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[ERR] %v", err)
		errHTTP400(ctx, err)
		return
	}

	argProd := db.CreateProductParams{
		Name:       req.Name,
		CategoryID: int32(req.CategoryID),
		ImageUrl:   req.Image,
		Price:      int32(req.Price),
		Stock:      int32(req.Stock),
		Sku:        utils.RandomString(),
	}

	prod, err := s.store.CreateProductWithReturn(context.Background(), argProd)
	if err != nil {
		log.Printf("[ERR] %v", err)
		errHTTP500(ctx)
		return
	}

	if req.Discount != nil {
		argDiscount := db.CreateDiscountParams{
			Qty:       int32(req.Discount.Qty),
			Type:      req.Discount.Type,
			ExpiredAt: time.Unix(req.Discount.ExpiredAt, 0),
			Result:    int32(req.Discount.Result),
		}
		argDiscount.ExpiredAtFormat = argDiscount.ExpiredAt.Format("02-01-2006")

		if req.Discount.Type == constants.PERCENT {
			discountedPrice := float64(prod.Price) - (float64(prod.Price) * float64(argDiscount.Result) / 100)
			argDiscount.StringFormat = fmt.Sprintf("Discount %d%% Rp. %.2f", argDiscount.Result, discountedPrice)
		} else if req.Discount.Type == constants.BUY_N {
			argDiscount.StringFormat = fmt.Sprintf("Buy %d only Rp. %d", argDiscount.Qty, argDiscount.Result)
		}

		discount, err := s.store.CreateDiscountWithReturn(context.Background(), argDiscount)
		if err != nil {
			log.Printf("[ERR] %v", err)
			errHTTP500(ctx)
			return
		}

		argProdDiscount := db.CreateProductDiscountParams{
			DiscountID: discount.ID,
			ProductID:  prod.ID,
		}
		if err := s.store.CreateProductDiscount(context.Background(), argProdDiscount); err != nil {
			log.Printf("[ERR] %v", err)
			errHTTP500(ctx)
			return
		}
	}

	ctx.JSON(http.StatusOK, dto.GenericResponse{
		Success: true,
		Message: "Success",
		Data: dto.Product{
			ID:         int(prod.ID),
			CategoryID: int(prod.CategoryID),
			SKU:        prod.Sku,
			Name:       prod.Name,
			Image:      prod.ImageUrl,
			Price:      int(prod.Price),
			Stock:      int(prod.Stock),
			CreatedAt:  &prod.CreatedAt,
			UpdatedAt:  &prod.UpdatedAt,
		},
	})
}

type getProductQuery struct {
	CategoryID int32  `form:"categoryId" binding:"omitempty"`
	Search     string `form:"q" binding:"omitempty"`
	PaginationQueryParams
}

// // todo: insert the category id and name into the response
func (s *Server) SearchProducts(ctx *gin.Context) {
	var query getProductQuery

	if err := ctx.ShouldBindQuery(&query); err != nil {
		log.Printf("[ERR] %v", err)
		errHTTP400(ctx, err)
		return
	}

	// don't accept both search and categoryID as query params
	if query.CategoryID != 0 && query.Search != "" {
		err := fmt.Errorf("Validation Failed on Query Parameters")
		log.Printf("[ERR] %v", err)
		errHTTP400(ctx, err)
		return
	}

	var products []dto.Product
	resp := dto.ListProducts{
		Meta: dto.Meta{
			Limit: int(query.Limit),
			Skip:  int(query.Skip),
		}}

	// filter with category id if present
	if query.CategoryID > 0 {
		arg := db.ListProductsDetailsByCategoryIDParams{
			CategoryID: query.CategoryID,
			Limit:      query.Limit,
			Offset:     query.Skip,
		}
		prods, err := s.store.ListProductsDetailsByCategoryID(context.Background(), arg)
		if err != nil {
			log.Printf("[ERR] %v", err)
			errHTTP500(ctx)
			return
		}

		resp.Meta.Total = len(prods)

		for i := 0; i < len(prods); i++ { //nolint
			p := dto.Product{
				ID:    int(prods[i].ID),
				SKU:   prods[i].Sku,
				Name:  prods[i].Name,
				Stock: int(prods[i].Stock),
				Price: int(prods[i].Price),
				Image: prods[i].ImageUrl,
				Category: &dto.Category{
					ID:   prods[i].CategoryID,
					Name: prods[i].CategoryName.String,
				},
			}
			if prods[i].DiscountID.Valid {
				p.Discount = &dto.Discount{
					ID:              prods[i].DiscountID.Int32,
					Qty:             prods[i].Qty.Int32,
					Type:            prods[i].Type.String,
					Result:          prods[i].Result.Int32,
					ExpiredAt:       prods[i].ExpiredAt.Time.String(),
					ExpiredAtFormat: prods[i].ExpiredAtFormat.String,
					StringFormat:    prods[i].StringFormat.String,
				}
			}
			products = append(products, p)
		}
	} else if query.Search != "" { // filter with search query if present
		sq := fmt.Sprintf("%%%s%%", query.Search)

		arg := db.ListProductsDetailsByNameParams{
			Name:   sq,
			Limit:  query.Limit,
			Offset: query.Skip,
		}
		prods, err := s.store.ListProductsDetailsByName(context.Background(), arg)
		if err != nil {
			log.Printf("[ERR] %v", err)
			errHTTP500(ctx)
			return
		}

		for i := 0; i < len(prods); i++ { // nolint
			p := dto.Product{
				ID:    int(prods[i].ID),
				SKU:   prods[i].Sku,
				Name:  prods[i].Name,
				Stock: int(prods[i].Stock),
				Price: int(prods[i].Price),
				Image: prods[i].ImageUrl,
				Category: &dto.Category{
					ID:   prods[i].CategoryID,
					Name: prods[i].CategoryName.String,
				},
			}
			if prods[i].DiscountID.Valid {
				p.Discount = &dto.Discount{
					ID:              prods[i].DiscountID.Int32,
					Qty:             prods[i].Qty.Int32,
					Type:            prods[i].Type.String,
					Result:          prods[i].Result.Int32,
					ExpiredAt:       prods[i].ExpiredAt.Time.String(),
					ExpiredAtFormat: prods[i].ExpiredAtFormat.String,
					StringFormat:    prods[i].StringFormat.String,
				}
			}
			products = append(products, p)
		}
	}

	resp.Product = products
	ctx.JSON(http.StatusOK, dto.GenericResponse{
		Success: true,
		Message: "Success",
		Data:    resp,
	})
}

// // todo: insert the category id and name into the response
func (s *Server) GetProduct(ctx *gin.Context) {
	var uri PathIDParam

	if err := ctx.ShouldBindUri(&uri); err != nil {
		log.Printf("[ERR] %v", err)
		errHTTP400(ctx, err)
		return
	}

	prod, err := s.store.GetProductDetails(context.Background(), uri.ID)
	if err != nil {
		log.Printf("[ERR] %v", err)
		if err == sql.ErrNoRows {
			errHTTP400(ctx, err)
			return
		}

		errHTTP500(ctx)
		return
	}

	resp := dto.Product{
		ID:    int(prod.ID),
		SKU:   prod.Sku,
		Name:  prod.Name,
		Stock: int(prod.Stock),
		Price: int(prod.Price),
		Image: prod.ImageUrl,
		Category: &dto.Category{
			ID:   prod.CategoryID,
			Name: prod.CategoryName,
		},
	}
	if prod.DiscountID.Valid {
		resp.Discount = &dto.Discount{
			ID:              prod.DiscountID.Int32,
			Qty:             prod.Qty.Int32,
			Type:            prod.Type.String,
			Result:          prod.Result.Int32,
			ExpiredAt:       prod.ExpiredAt.Time.String(),
			ExpiredAtFormat: prod.ExpiredAtFormat.String,
			StringFormat:    prod.StringFormat.String,
		}
	}

	ctx.JSON(http.StatusOK, dto.GenericResponse{
		Success: true,
		Message: "Success",
		Data:    resp,
	})
}

type updateProductBody struct {
	CategoryID int32  `json:"categoryId"`
	Name       string `json:"name"`
	Image      string `json:"image"`
	Price      int32  `json:"price"`
	Stock      int32  `json:"stock"`
}

func (s *Server) UpdateProduct(ctx *gin.Context) {
	var uri PathIDParam
	var body updateProductBody

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

	prod, err := s.store.GetProduct(context.Background(), uri.ID)
	if err != nil {
		log.Printf("[ERR] %v", err)
		if err == sql.ErrNoRows {
			errHTTP404(ctx)
			return
		}

		errHTTP500(ctx)
		return
	}

	if body.CategoryID == 0 {
		body.CategoryID = prod.CategoryID
	}
	if body.Name == "" {
		body.Name = prod.Name
	}
	if body.Image == "" {
		body.Image = prod.ImageUrl
	}
	if body.Price == 0 {
		body.Price = prod.Price
	}
	if body.Stock == 0 {
		body.Stock = prod.Stock
	}

	arg := db.UpdateProductParams{
		CategoryID: body.CategoryID,
		Name:       body.Name,
		ImageUrl:   body.Image,
		Price:      body.Price,
		Stock:      body.Stock,
		ID:         uri.ID,
	}
	if err := s.store.UpdateProduct(context.Background(), arg); err != nil {
		errHTTP500(ctx)
		return
	}

	ctx.JSON(http.StatusOK, dto.GenericResponse{
		Success: true,
		Message: "Success",
	})
}

func (s *Server) DeleteProduct(ctx *gin.Context) { //nolint
	var uri PathIDParam

	if err := ctx.ShouldBindUri(&uri); err != nil {
		log.Printf("[ERR] %v", err)
		errHTTP400(ctx, err)
		return
	}

	_, err := s.store.GetProduct(context.Background(), uri.ID)
	if err != nil {
		log.Printf("[ERR] %v", err)
		if err == sql.ErrNoRows {
			errHTTP400(ctx, err)
			return
		}

		errHTTP500(ctx)
		return
	}
	if err := s.store.DeleteProduct(context.Background(), uri.ID); err != nil {
		errHTTP500(ctx)
		return
	}

	ctx.JSON(http.StatusOK, dto.GenericResponse{
		Success: true,
		Message: "Success",
	})
}
