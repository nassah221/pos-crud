package dto

import "time"

type Payment struct {
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	Name      string     `json:"name"`
	Type      string     `json:"type"`
	Logo      string     `json:"logo,omitempty"`
	ID        int        `json:"paymentId"`
}

type ListPayment struct {
	Payment []Payment `json:"payments"`
	Meta    Meta      `json:"meta"`
}
