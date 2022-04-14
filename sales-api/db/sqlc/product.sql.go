// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0
// source: product.sql

package db

import (
	"context"
	"database/sql"
)

const createProduct = `-- name: CreateProduct :execresult
INSERT INTO products (
    name,image_url,price,stock,category_id,sku
) VALUES (
    ?,?,?,?,?,?
)
`

type CreateProductParams struct {
	Name       string `json:"name"`
	ImageUrl   string `json:"image_url"`
	Price      int32  `json:"price"`
	Stock      int32  `json:"stock"`
	CategoryID int32  `json:"category_id"`
	Sku        string `json:"sku"`
}

func (q *Queries) CreateProduct(ctx context.Context, arg CreateProductParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, createProduct,
		arg.Name,
		arg.ImageUrl,
		arg.Price,
		arg.Stock,
		arg.CategoryID,
		arg.Sku,
	)
}

const deleteProduct = `-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = ?
`

func (q *Queries) DeleteProduct(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteProduct, id)
	return err
}

const getProduct = `-- name: GetProduct :one
SELECT id, sku, name, stock, price, image_url, category_id, created_at, updated_at FROM products
WHERE id = ? LIMIT 1
`

func (q *Queries) GetProduct(ctx context.Context, id int32) (Product, error) {
	row := q.db.QueryRowContext(ctx, getProduct, id)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Sku,
		&i.Name,
		&i.Stock,
		&i.Price,
		&i.ImageUrl,
		&i.CategoryID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getProductDetails = `-- name: GetProductDetails :one
SELECT p.name,p.id,p.sku,p.stock,p.price,p.image_url,p.category_id,
dis.id AS discount_id,dis.qty,dis.type,dis.expired_at,dis.result,dis.expired_at_format,dis.string_format,
c.name AS category_name FROM products AS p JOIN categories AS c ON c.id = p.category_id 
LEFT JOIN product_discount AS d ON p.id = d.product_id LEFT JOIN discounts AS dis ON dis.id=discount_id 
WHERE p.id = ?
`

type GetProductDetailsRow struct {
	Name            string         `json:"name"`
	ID              int32          `json:"id"`
	Sku             string         `json:"sku"`
	Stock           int32          `json:"stock"`
	Price           int32          `json:"price"`
	ImageUrl        string         `json:"image_url"`
	CategoryID      int32          `json:"category_id"`
	DiscountID      sql.NullInt32  `json:"discount_id"`
	Qty             sql.NullInt32  `json:"qty"`
	Type            sql.NullString `json:"type"`
	ExpiredAt       sql.NullTime   `json:"expired_at"`
	Result          sql.NullInt32  `json:"result"`
	ExpiredAtFormat sql.NullTime   `json:"expired_at_format"`
	StringFormat    sql.NullString `json:"string_format"`
	CategoryName    string         `json:"category_name"`
}

func (q *Queries) GetProductDetails(ctx context.Context, id int32) (GetProductDetailsRow, error) {
	row := q.db.QueryRowContext(ctx, getProductDetails, id)
	var i GetProductDetailsRow
	err := row.Scan(
		&i.Name,
		&i.ID,
		&i.Sku,
		&i.Stock,
		&i.Price,
		&i.ImageUrl,
		&i.CategoryID,
		&i.DiscountID,
		&i.Qty,
		&i.Type,
		&i.ExpiredAt,
		&i.Result,
		&i.ExpiredAtFormat,
		&i.StringFormat,
		&i.CategoryName,
	)
	return i, err
}

const listProducts = `-- name: ListProducts :many
SELECT id, sku, name, stock, price, image_url, category_id, created_at, updated_at FROM products
ORDER BY id LIMIT ? OFFSET ?
`

type ListProductsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListProducts(ctx context.Context, arg ListProductsParams) ([]Product, error) {
	rows, err := q.db.QueryContext(ctx, listProducts, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Product
	for rows.Next() {
		var i Product
		if err := rows.Scan(
			&i.ID,
			&i.Sku,
			&i.Name,
			&i.Stock,
			&i.Price,
			&i.ImageUrl,
			&i.CategoryID,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listProductsDetailsByCategoryID = `-- name: ListProductsDetailsByCategoryID :many
select p.id,p.name,p.sku,p.image_url,p.stock,p.price,p.category_id,
c.name as category_name,
dis.id as discount_id,
d.qty,d.type,d.result,d.expired_at,d.expired_at_format,d.string_format
from products as p left join categories as c on c.id = p.category_id
left join product_discount as dis on dis.product_id = p.id
left join discounts as d on d.id=discount_id
where category_id=? ORDER BY p.id LIMIT ? OFFSET ?
`

type ListProductsDetailsByCategoryIDParams struct {
	CategoryID int32 `json:"category_id"`
	Limit      int32 `json:"limit"`
	Offset     int32 `json:"offset"`
}

type ListProductsDetailsByCategoryIDRow struct {
	ID              int32          `json:"id"`
	Name            string         `json:"name"`
	Sku             string         `json:"sku"`
	ImageUrl        string         `json:"image_url"`
	Stock           int32          `json:"stock"`
	Price           int32          `json:"price"`
	CategoryID      int32          `json:"category_id"`
	CategoryName    sql.NullString `json:"category_name"`
	DiscountID      sql.NullInt32  `json:"discount_id"`
	Qty             sql.NullInt32  `json:"qty"`
	Type            sql.NullString `json:"type"`
	Result          sql.NullInt32  `json:"result"`
	ExpiredAt       sql.NullTime   `json:"expired_at"`
	ExpiredAtFormat sql.NullTime   `json:"expired_at_format"`
	StringFormat    sql.NullString `json:"string_format"`
}

func (q *Queries) ListProductsDetailsByCategoryID(ctx context.Context, arg ListProductsDetailsByCategoryIDParams) ([]ListProductsDetailsByCategoryIDRow, error) {
	rows, err := q.db.QueryContext(ctx, listProductsDetailsByCategoryID, arg.CategoryID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListProductsDetailsByCategoryIDRow
	for rows.Next() {
		var i ListProductsDetailsByCategoryIDRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Sku,
			&i.ImageUrl,
			&i.Stock,
			&i.Price,
			&i.CategoryID,
			&i.CategoryName,
			&i.DiscountID,
			&i.Qty,
			&i.Type,
			&i.Result,
			&i.ExpiredAt,
			&i.ExpiredAtFormat,
			&i.StringFormat,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listProductsDetailsByName = `-- name: ListProductsDetailsByName :many
select p.id,p.name,p.sku,p.image_url,p.stock,p.price,p.category_id,
c.name as category_name,
dis.id as discount_id,
d.qty,d.type,d.result,d.expired_at,d.expired_at_format,d.string_format
from products as p left join categories as c on c.id = p.category_id
left join product_discount as dis on dis.product_id = p.id
left join discounts as d on d.id=discount_id
WHERE p.name LIKE ? ORDER BY p.id LIMIT ? OFFSET ?
`

type ListProductsDetailsByNameParams struct {
	Name   string `json:"name"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}

type ListProductsDetailsByNameRow struct {
	ID              int32          `json:"id"`
	Name            string         `json:"name"`
	Sku             string         `json:"sku"`
	ImageUrl        string         `json:"image_url"`
	Stock           int32          `json:"stock"`
	Price           int32          `json:"price"`
	CategoryID      int32          `json:"category_id"`
	CategoryName    sql.NullString `json:"category_name"`
	DiscountID      sql.NullInt32  `json:"discount_id"`
	Qty             sql.NullInt32  `json:"qty"`
	Type            sql.NullString `json:"type"`
	Result          sql.NullInt32  `json:"result"`
	ExpiredAt       sql.NullTime   `json:"expired_at"`
	ExpiredAtFormat sql.NullTime   `json:"expired_at_format"`
	StringFormat    sql.NullString `json:"string_format"`
}

func (q *Queries) ListProductsDetailsByName(ctx context.Context, arg ListProductsDetailsByNameParams) ([]ListProductsDetailsByNameRow, error) {
	rows, err := q.db.QueryContext(ctx, listProductsDetailsByName, arg.Name, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListProductsDetailsByNameRow
	for rows.Next() {
		var i ListProductsDetailsByNameRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Sku,
			&i.ImageUrl,
			&i.Stock,
			&i.Price,
			&i.CategoryID,
			&i.CategoryName,
			&i.DiscountID,
			&i.Qty,
			&i.Type,
			&i.Result,
			&i.ExpiredAt,
			&i.ExpiredAtFormat,
			&i.StringFormat,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateProduct = `-- name: UpdateProduct :exec
UPDATE products SET category_id=?,name=?,image_url=?,price=?,stock=?,updated_at=CURRENT_TIMESTAMP
WHERE id = ?
`

type UpdateProductParams struct {
	CategoryID int32  `json:"category_id"`
	Name       string `json:"name"`
	ImageUrl   string `json:"image_url"`
	Price      int32  `json:"price"`
	Stock      int32  `json:"stock"`
	ID         int32  `json:"id"`
}

func (q *Queries) UpdateProduct(ctx context.Context, arg UpdateProductParams) error {
	_, err := q.db.ExecContext(ctx, updateProduct,
		arg.CategoryID,
		arg.Name,
		arg.ImageUrl,
		arg.Price,
		arg.Stock,
		arg.ID,
	)
	return err
}
