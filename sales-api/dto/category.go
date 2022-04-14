package dto

import "time"

type Category struct {
	ID        int32      `json:"categoryId"`
	Name      string     `json:"name"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
}

type ListCategory struct {
	Category []Category `json:"categories"`
	Meta     Meta       `json:"meta"`
}
