package dto

import "time"

type Cashier struct {
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	Name      string     `json:"name"`
	Passcode  string     `json:"passcode,omitempty"`
	ID        int        `json:"cashierId"`
}

type ListCashiers struct {
	Cashiers []Cashier `json:"cashiers"`
	Meta     Meta      `json:"meta"`
}
