package dto

import "time"

type ProductDiscountRequest struct {
	Qty       int    `json:"qty" binding:"required,min=1"`
	Type      string `json:"type" binding:"required,oneof=BUY_N PERCENT"`
	Result    int    `json:"result" binding:"required,min=1"`
	ExpiredAt int64  `json:"expiredAt" binding:"required"`
}

type ProductCreateRequest struct {
	CategoryID int                     `json:"categoryId" binding:"required"`
	Name       string                  `json:"name" binding:"required"`
	Image      string                  `json:"image" binding:"required"`
	Price      int                     `json:"price" binding:"required"`
	Stock      int                     `json:"stock" binding:"required"`
	Discount   *ProductDiscountRequest `json:"discount,omitempty"`
}

type Product struct {
	ID         int        `json:"productId"`
	Name       string     `json:"name"`
	Image      string     `json:"image"`
	Price      int        `json:"price"`
	Stock      int        `json:"stock"`
	SKU        string     `json:"sku"`
	CategoryID int        `json:"categoriesId,omitempty"`
	Discount   *Discount  `json:"discount,omitempty"`
	Category   *Category  `json:"category,omitempty"`
	CreatedAt  *time.Time `json:"createdAt,omitempty"`
	UpdatedAt  *time.Time `json:"updatedAt,omitempty"`
}

type ListProducts struct {
	Product []Product `json:"products"`
	Meta    Meta      `json:"meta,omitempty"`
}
