package dto

type Discount struct {
	StringFormat    string `json:"stringFormat,omitempty"`
	Type            string `json:"type"`
	ExpiredAt       string `json:"expiredAt"`
	ExpiredAtFormat string `json:"expiredAtFormat,omitempty"`
	Qty             int32  `json:"qty"`
	Result          int32  `json:"result"`
	ID              int32  `json:"discountId"`
}
