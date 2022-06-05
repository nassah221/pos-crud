package dto

import "time"

type Category struct {
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	Name      string     `json:"name"`
	ID        int32      `json:"categoryId"`
}

type ListCategory struct {
	Category []Category `json:"categories"`
	Meta     Meta       `json:"meta"`
}
