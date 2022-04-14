package dto

import "time"

type OrderProductRequest struct {
	ProductID int32 `json:"productId" binding:"required"`
	Qty       int   `json:"qty" binding:"required"`
}

type OrderAddRequest struct {
	PaymentID int                   `json:"paymentId" binding:"required"`
	TotalPaid int                   `json:"totalPaid" binding:"required"`
	Products  []OrderProductRequest `json:"products" binding:"required"`
}

type OrderProduct struct {
	ProductID        int       `json:"productId"`
	Name             string    `json:"name"`
	Price            int       `json:"price"`
	DiscountID       int32     `json:"discountsId,omitempty"`
	Discount         *Discount `json:"discount,omitempty"`
	Qty              int       `json:"qty"`
	TotalNormalPrice float64   `json:"totalNormalPrice"`
	TotalFinalPrice  float64   `json:"totalFinalPrice"`
}

type OrderAdd struct {
	ID             int            `json:"orderId"`
	CashierID      int            `json:"cashiersId"`
	PaymentTypesID int            `json:"paymentTypesId"`
	TotalPrice     float64        `json:"totalPrice"`
	TotalPaid      int            `json:"totalPaid"`
	TotalReturn    float64        `json:"totalReturn"`
	ReceiptID      string         `json:"receiptId"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	Prodcts        []OrderProduct `json:"products"`
}

type OrderDetails struct {
	OrderID       int32          `json:"orderId"`
	CashierID     int32          `json:"cashiersId"`
	PaymentTypeID int32          `json:"paymentTypesId"`
	TotalPrice    float64        `json:"totalPrice"`
	TotalPaid     float64        `json:"totalPaid"`
	TotalReturn   float64        `json:"totalReturn"`
	ReceiptID     string         `json:"receiptId"`
	CreatedAt     time.Time      `json:"createdAt"`
	Cashier       Cashier        `json:"cashier"`
	PaymentType   Payment        `json:"payment_type"`
	Products      []OrderProduct `json:"products"`
}

type ListAllOrderDetails struct {
	Orders []OrderDetails `json:"orders"`
	Meta   Meta           `json:"meta"`
}
