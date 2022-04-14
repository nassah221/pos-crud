package dto

type Discount struct {
	ID              int32  `json:"discountId"`
	Qty             int32  `json:"qty"`
	Type            string `json:"type"`
	Result          int32  `json:"result"`
	ExpiredAt       string `json:"expiredAt"` // can be in unix nano and UTC
	ExpiredAtFormat string `json:"expiredAtFormat,omitempty"`
	StringFormat    string `json:"stringFormat,omitempty"`
}
