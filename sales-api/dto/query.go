package dto

type RevenuePayment struct {
	ID          int32   `json:"paymentTypeId"`
	Name        string  `json:"name"`
	Logo        string  `json:"logo,omitempty"`
	TotalAmount float64 `json:"totalAmount"`
}

type Revenue struct {
	TotalRevenue float64          `json:"totalRevenue"`
	PaymentType  []RevenuePayment `json:"paymentTypes"`
}
