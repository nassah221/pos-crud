package dto

import "time"

type OrderProductRequest struct {
	ProductID int32 `json:"productId" binding:"required"`
	Qty       int   `json:"qty" binding:"required"`
}

type OrderAddRequest struct {
	Products  []OrderProductRequest `json:"products" binding:"required"`
	PaymentID int                   `json:"paymentId" binding:"required"`
	TotalPaid int                   `json:"totalPaid" binding:"required"`
}

type OrderProduct struct {
	Discount         *Discount `json:"discount,omitempty"`
	Name             string    `json:"name"`
	ProductID        int       `json:"productId"`
	Price            int       `json:"price"`
	Qty              int       `json:"qty"`
	TotalNormalPrice float64   `json:"totalNormalPrice"`
	TotalFinalPrice  float64   `json:"totalFinalPrice"`
	DiscountID       int32     `json:"discountsId,omitempty"`
}

type OrderAdd struct {
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	ReceiptID      string         `json:"receiptId"`
	Prodcts        []OrderProduct `json:"products"`
	TotalPaid      int            `json:"totalPaid"`
	TotalReturn    float64        `json:"totalReturn"`
	ID             int            `json:"orderId"`
	CashierID      int            `json:"cashiersId"`
	PaymentTypesID int            `json:"paymentTypesId"`
	TotalPrice     float64        `json:"totalPrice"`
}

type OrderDetails struct {
	CreatedAt     time.Time      `json:"createdAt"`
	PaymentType   Payment        `json:"payment_type"`
	ReceiptID     string         `json:"receiptId"`
	Cashier       Cashier        `json:"cashier"`
	Products      []OrderProduct `json:"products"`
	TotalReturn   float64        `json:"totalReturn"`
	TotalPrice    float64        `json:"totalPrice"`
	TotalPaid     float64        `json:"totalPaid"`
	OrderID       int32          `json:"orderId"`
	CashierID     int32          `json:"cashiersId"`
	PaymentTypeID int32          `json:"paymentTypesId"`
}

type ListAllOrderDetails struct {
	Orders []OrderDetails `json:"orders"`
	Meta   Meta           `json:"meta"`
}
