package api

import (
	"net/http"
	db "sales-api/db/sqlc"
	"sales-api/dto"
	"sales-api/errors"

	"github.com/gin-gonic/gin"
)

func errHTTP500(ctx *gin.Context) {
	ctx.JSON(http.StatusInternalServerError, dto.GenericResponse{
		Success: false,
		Error:   ErrInternalServer,
	})
}

func errHTTP400(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusBadRequest, dto.GenericResponse{
		Success: false,
		Message: MsgInvalidRequest,
		Error:   err.Error(),
	})
}

func errHTTP400BodyInvalid(ctx *gin.Context, msg string, err []errors.BodyValidationError) {
	ctx.JSON(http.StatusBadRequest, dto.GenericResponse{
		Success: false,
		Message: msg,
		Error:   err,
	})
}

func errHTTP404(ctx *gin.Context) {
	ctx.JSON(http.StatusNotFound, dto.GenericResponse{
		Success: false,
		Message: MsgNotFound,
	})
}

func populateOrderDetails(o []db.GetOrderDetailsRow) dto.OrderDetails {
	op := []dto.OrderProduct{}
	od := dto.OrderDetails{
		OrderID:       o[0].Orderid,
		CashierID:     o[0].Cashiersid,
		PaymentTypeID: o[0].Paymenttypesid,
		TotalPrice:    float64(o[0].TotalPrice),
		TotalPaid:     float64(o[0].PricePaid),
		TotalReturn:   float64(o[0].TotalReturn),
		CreatedAt:     o[0].CreatedAt,
		ReceiptID:     o[0].ReceiptID,
		Cashier: dto.Cashier{
			ID:   int(o[0].Cashiersid),
			Name: o[0].Cashiername,
		},
		PaymentType: dto.Payment{
			ID:   int(o[0].Paymenttypeid),
			Name: o[0].Paymentname,
			Logo: o[0].Logo.String,
			Type: o[0].Paymenttype,
		},
	}
	for i := 0; i < len(o); i++ {
		prod := dto.OrderProduct{
			ProductID:        int(o[i].ProductID),
			Name:             o[i].ProductName,
			DiscountID:       o[i].Discountsid.Int32,
			Qty:              int(o[i].Productqty),
			Price:            int(o[i].Price),
			TotalFinalPrice:  float64(o[i].TotalFinalPrice),
			TotalNormalPrice: float64(o[i].TotalNormalPrice),
		}
		if o[i].Discountsid.Valid {
			prod.Discount = &dto.Discount{
				ID:              o[i].Discountsid.Int32,
				Qty:             o[i].Discountqty.Int32,
				Type:            o[i].Discounttype.String,
				Result:          o[i].Result.Int32,
				ExpiredAt:       o[i].ExpiredAt.Time.String(),
				ExpiredAtFormat: o[i].ExpiredAtFormat.String,
				StringFormat:    o[i].StringFormat.String,
			}
		}
		op = append(op, prod)
	}
	od.Products = op
	return od
}
