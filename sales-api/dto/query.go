package dto

type RevenuePayment struct {
	Name        string  `json:"name"`
	Logo        string  `json:"logo,omitempty"`
	TotalAmount float64 `json:"totalAmount"`
	ID          int32   `json:"paymentTypeId"`
}

type Revenue struct {
	PaymentType  []RevenuePayment `json:"paymentTypes"`
	TotalRevenue float64          `json:"totalRevenue"`
}
