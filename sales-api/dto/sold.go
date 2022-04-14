package dto

type Sold struct {
	ProductID   int32  `json:"productId"`
	ProductName string `json:"productName"`
	TotalQty    int32  `json:"totalQty"`
	TotalAmount int64  `json:"totalAmount"`
}

type SoldResponse struct {
	OrderProducts []Sold `json:"orderProducts"`
}
