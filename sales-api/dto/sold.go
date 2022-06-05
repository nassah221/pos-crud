package dto

type Sold struct {
	ProductName string `json:"productName"`
	TotalAmount int64  `json:"totalAmount"`
	ProductID   int32  `json:"productId"`
	TotalQty    int32  `json:"totalQty"`
}

type SoldResponse struct {
	OrderProducts []Sold `json:"orderProducts"`
}
