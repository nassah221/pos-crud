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
	"sales-api/errors"
	"sales-api/token"
	"sales-api/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *Server) ListAllOrderDetails(ctx *gin.Context) {
	var query PaginationQueryParams

	if err := ctx.ShouldBindQuery(&query); err != nil {
		log.Printf("[ERR] %v", err)
		errHTTP400(ctx, err)
		return
	}

	allOrderDetails := []dto.OrderDetails{}
	for i := query.Skip + 9; i <= query.Skip+query.Limit+9; i++ {
		o, err := s.store.GetOrderDetails(context.Background(), i)
		if err != nil {
			log.Printf("[ERR] %v", err)
			errHTTP500(ctx)
			return
		}
		if len(o) == 0 {
			break
		}
		od := populateOrderDetails(o)
		allOrderDetails = append(allOrderDetails, od)
	}

	ctx.JSON(http.StatusOK, dto.GenericResponse{
		Success: true,
		Message: "Success",
		Data: dto.ListAllOrderDetails{
			Orders: allOrderDetails,
			Meta: dto.Meta{
				Limit: int(query.Limit),
				Skip:  int(query.Skip),
				Total: len(allOrderDetails),
			},
		},
	})
}

func (s *Server) GetOrderDetails(ctx *gin.Context) {
	var uri PathIDParam

	if err := ctx.ShouldBindUri(&uri); err != nil {
		log.Printf("[ERR] %v", err)
		errHTTP400(ctx, err)
		return
	}

	o, err := s.store.GetOrderDetails(context.Background(), uri.ID)
	if err != nil {
		log.Printf("[ERR] %v", err)
		errHTTP500(ctx)
		return
	}
	if len(o) == 0 {
		log.Printf("[INFO] order id %d not found ", uri.ID)
		errHTTP404(ctx, constants.Order)
		return
	}
	od := populateOrderDetails(o)

	ctx.JSON(http.StatusOK, dto.GenericResponse{
		Success: true,
		Message: "Success",
		Data: map[string]interface{}{
			"order": od,
		},
	})
}

func (s *Server) CreateOrder(ctx *gin.Context) {
	var body dto.OrderAddRequest

	if err := ctx.ShouldBindJSON(&body); err != nil {
		log.Printf("[ERR] %v", err)
		vErr, msg := errors.FromFieldValidationErrorPOST(err)
		errHTTP400BodyInvalid(ctx, msg, vErr)
		return
	}

	payload := ctx.MustGet(authPayloadKey)
	cashier, _ := payload.(*token.Payload) //nolint

	prods, receiptTotal, err := s.orderTotal(body.Products, false)
	if err != nil {
		errHTTP500(ctx)
		return
	}

	totalReturn := float64(body.TotalPaid) - receiptTotal
	receiptID := utils.RandomString()

	arg := db.CreateOrderParams{
		CashierID:   cashier.CashierID,
		PaymentID:   int32(body.PaymentID),
		TotalPrice:  int32(receiptTotal),
		PricePaid:   int32(body.TotalPaid),
		TotalReturn: int32(totalReturn),
		ReceiptID:   receiptID,
	}
	res, err := s.store.CreateOrder(context.Background(), arg)
	if err != nil {
		log.Printf("[ERR] %v", err)
		errHTTP500(ctx)
		return
	}

	insertID, err := res.LastInsertId()
	if err != nil {
		insertErr := fmt.Errorf("cannot fetch last insert id: %v", err)
		log.Printf("[ERR] %v", insertErr)
		errHTTP500(ctx)
	}

	if err := s.insertOrderDetails(prods, int32(insertID)); err != nil {
		log.Printf("[ERR] %v", err)
		errHTTP500(ctx)
	}

	resp := dto.OrderAdd{
		ID:             int(insertID),
		CashierID:      int(cashier.CashierID),
		PaymentTypesID: body.PaymentID,
		TotalPrice:     receiptTotal,
		TotalPaid:      body.TotalPaid,
		TotalReturn:    totalReturn,
		ReceiptID:      receiptID,
		UpdatedAt:      time.Now(),
		CreatedAt:      time.Now(),
		Prodcts:        prods,
	}

	ctx.JSON(http.StatusOK, dto.GenericResponse{
		Success: true,
		Message: "Success",
		Data:    resp,
	})
}

func (s *Server) SubtotalOrder(ctx *gin.Context) {
	var body []dto.OrderProductRequest

	if err := ctx.ShouldBindJSON(&body); err != nil {
		log.Printf("[ERR] %v", err)
		// errHTTP400(ctx, err)
		ctx.JSON(http.StatusBadRequest, dto.GenericResponse{
			Success: false,
			Message: "body ValidationError: \"value\" must be an array",
			Error: []map[string]interface{}{
				{
					"message": "\"value\" must be an array",
					"path":    []string{},
					"type":    "array.base",
					"context": map[string]interface{}{
						"label": "value",
						"value": struct{}{},
					},
				},
			},
		})
		return
	}

	prods, receiptTotal, err := s.orderTotal(body, true)
	if err != nil {
		errHTTP500(ctx)
		return
	}

	resp := map[string]interface{}{
		"subtotal": receiptTotal,
		"products": prods,
	}

	ctx.JSON(http.StatusOK, dto.GenericResponse{
		Success: true,
		Message: "Success",
		Data:    resp,
	})
}

func (s *Server) insertOrderDetails(body []dto.OrderProduct, orderID int32) error {
	for i := 0; i < len(body); i++ {
		id := (body[i].ProductID)
		prod, err := s.store.GetProductDetails(context.Background(), int32(id))
		if err != nil {
			return err
		}
		orderArg := db.CreateOrderDetailsParams{
			OrderID:          orderID,
			ProductID:        prod.ID,
			ProductName:      prod.Name,
			Qty:              int32(body[i].Qty),
			Price:            prod.Price,
			TotalNormalPrice: int64(body[i].TotalNormalPrice),
			TotalFinalPrice:  int64(body[i].TotalFinalPrice),
		}

		if !prod.DiscountID.Valid {
			orderArg.DiscountID = sql.NullInt32{}
		}
		orderArg.DiscountID = prod.DiscountID

		// order_details
		if err := s.store.CreateOrderDetails(context.Background(), orderArg); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) orderTotal(body []dto.OrderProductRequest, calculateOnly bool) ([]dto.OrderProduct, float64, error) {
	var receiptTotal float64
	var prods []dto.OrderProduct
	for i := 0; i < len(body); i++ {
		id := body[i].ProductID
		prod, err := s.store.GetProductDetails(context.Background(), id)
		if err != nil {
			return nil, 0, err
		}
		// If the product doesn't have a discount
		if !prod.DiscountID.Valid {
			p := dto.OrderProduct{
				ProductID: int(prod.ID),
				Name:      prod.Name,
				Price:     int(prod.Price),
				Qty:       body[i].Qty,
			}
			p.TotalNormalPrice = float64(body[i].Qty * p.Price)
			p.TotalFinalPrice = p.TotalNormalPrice

			// Sum all products price
			receiptTotal += p.TotalFinalPrice
			prods = append(prods, p)

			if !calculateOnly {
				// update the product stock
				arg := db.UpdateProductParams{
					Name:       prod.Name,
					ID:         id,
					CategoryID: prod.CategoryID,
					Stock:      prod.Stock - int32(body[i].Qty),
					Price:      prod.Price,
					ImageUrl:   prod.ImageUrl,
				}
				if err := s.store.UpdateProduct(context.Background(), arg); err != nil {
					log.Printf("[ERR] %v", err)
					return nil, 0, err
				}

				// update sold items
				soldArg := db.InsertSoldParams{
					ProductID:   id,
					ProductName: prod.Name,
					TotalQty:    int32(body[i].Qty),
					TotalAmount: int64(p.TotalFinalPrice),
				}
				if err := s.store.InsertSold(context.Background(), soldArg); err != nil {
					log.Printf("[ERR] %v", err)
					return nil, 0, err
				}
			}
			continue
		}

		p := dto.OrderProduct{
			ProductID: int(prod.ID),
			Name:      prod.Name,
			Price:     int(prod.Price),
			Qty:       body[i].Qty,
			Discount: &dto.Discount{
				ID:              prod.DiscountID.Int32,
				Qty:             prod.Qty.Int32,
				Type:            prod.Type.String,
				Result:          prod.Result.Int32,
				ExpiredAt:       prod.ExpiredAt.Time.String(),
				ExpiredAtFormat: prod.ExpiredAt.Time.String(),
				StringFormat:    prod.StringFormat.String,
			},
		}
		p.TotalNormalPrice = float64(body[i].Qty * p.Price)

		switch prod.Type.String {
		case "BUY_N":
			if body[i].Qty < int(prod.Qty.Int32) {
				// if discount quantity is greater than order quantity then discount is not applied
				p.TotalFinalPrice = p.TotalNormalPrice
			} else {
				remainderQty := body[i].Qty - int(prod.Qty.Int32)
				p.TotalFinalPrice = float64(remainderQty*p.Price) + float64(prod.Result.Int32)
			}
		case "PERCENT":
			if body[i].Qty < int(prod.Qty.Int32) {
				// if discount quantity is greater than order quantity then discount is not applied
				p.TotalFinalPrice = p.TotalNormalPrice
			} else {
				// If the discount percentage is per single product
				result := float64(prod.Result.Int32) / 100
				if prod.Qty.Int32 == 1 {
					discount := result * float64(p.Price) * float64(p.Qty)
					p.TotalFinalPrice = p.TotalNormalPrice - discount
				} else { // If the discount percent is per more than one product
					discount := result * float64(p.Price) * float64(prod.Qty.Int32)
					p.TotalFinalPrice = p.TotalNormalPrice - discount
				}
			}
		}
		if !calculateOnly {
			// update the stock
			arg := db.UpdateProductParams{
				ID:         id,
				Name:       prod.Name,
				CategoryID: prod.CategoryID,
				Stock:      prod.Stock - int32(body[i].Qty),
				Price:      prod.Price,
				ImageUrl:   prod.ImageUrl,
			}
			if err := s.store.UpdateProduct(context.Background(), arg); err != nil {
				log.Printf("[ERR] %v", err)
				return nil, 0, err
			}

			// update sold items
			soldArg := db.InsertSoldParams{
				ProductID:   id,
				ProductName: prod.Name,
				TotalQty:    int32(body[i].Qty),
				TotalAmount: int64(p.TotalFinalPrice),
			}
			if err := s.store.InsertSold(context.Background(), soldArg); err != nil {
				log.Printf("[ERR] %v", err)
				return nil, 0, err
			}
		}

		receiptTotal += p.TotalFinalPrice
		prods = append(prods, p)
	}
	return prods, receiptTotal, nil
}
