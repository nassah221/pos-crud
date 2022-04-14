package dto

type Meta struct {
	Limit int `json:"limit"`
	Skip  int `json:"skip"`
	Total int `json:"total"`
}
