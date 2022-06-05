package dto

import "time"

type ProductDiscountRequest struct {
	Type      string `json:"type" binding:"required,oneof=BUY_N PERCENT"`
	Qty       int    `json:"qty" binding:"required,min=1"`
	Result    int    `json:"result" binding:"required,min=1"`
	ExpiredAt int64  `json:"expiredAt" binding:"required"`
}

type ProductCreateRequest struct {
	Discount   *ProductDiscountRequest `json:"discount,omitempty"`
	Name       string                  `json:"name" binding:"required"`
	Image      string                  `json:"image" binding:"required"`
	CategoryID int                     `json:"categoryId" binding:"required"`
	Price      int                     `json:"price" binding:"required"`
	Stock      int                     `json:"stock" binding:"required"`
}

type Product struct {
	Discount   *Discount  `json:"discount,omitempty"`
	Category   *Category  `json:"category,omitempty"`
	CreatedAt  *time.Time `json:"createdAt,omitempty"`
	UpdatedAt  *time.Time `json:"updatedAt,omitempty"`
	SKU        string     `json:"sku"`
	Name       string     `json:"name"`
	Image      string     `json:"image"`
	ID         int        `json:"productId"`
	CategoryID int        `json:"categoriesId,omitempty"`
	Price      int        `json:"price"`
	Stock      int        `json:"stock"`
}

type ListProducts struct {
	Product []Product `json:"products"`
	Meta    Meta      `json:"meta,omitempty"`
}
