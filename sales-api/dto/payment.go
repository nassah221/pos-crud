package dto

import "time"

type Payment struct {
	ID        int        `json:"paymentId"`
	Name      string     `json:"name"`
	Type      string     `json:"type"`
	Logo      string     `json:"logo,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
}

type ListPayment struct {
	Payment []Payment `json:"payments"`
	Meta    Meta      `json:"meta"`
}
