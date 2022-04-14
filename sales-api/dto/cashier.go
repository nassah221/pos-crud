package dto

import "time"

type Cashier struct {
	ID        int        `json:"cashierId"`
	Name      string     `json:"name"`
	Passcode  string     `json:"passcode,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
}

type ListCashiers struct {
	Cashiers []Cashier `json:"cashiers"`
	Meta     Meta      `json:"meta"`
}
