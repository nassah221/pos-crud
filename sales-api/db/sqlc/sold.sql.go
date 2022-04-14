// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0
// source: sold.sql

package db

import (
	"context"
)

const getSold = `-- name: GetSold :many
SELECT id,product_id,product_name,SUM(total_qty) AS total_qty,SUM(total_amount) AS total_Amount 
FROM sold GROUP BY product_id ORDER BY id
`

type GetSoldRow struct {
	ID          int32       `json:"id"`
	ProductID   int32       `json:"product_id"`
	ProductName string      `json:"product_name"`
	TotalQty    interface{} `json:"total_qty"`
	TotalAmount interface{} `json:"total_amount"`
}

func (q *Queries) GetSold(ctx context.Context) ([]GetSoldRow, error) {
	rows, err := q.db.QueryContext(ctx, getSold)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetSoldRow
	for rows.Next() {
		var i GetSoldRow
		if err := rows.Scan(
			&i.ID,
			&i.ProductID,
			&i.ProductName,
			&i.TotalQty,
			&i.TotalAmount,
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

const insertSold = `-- name: InsertSold :exec
INSERT INTO sold (
    product_id, product_name, total_qty, total_amount
) VALUES (
    ?,?,?,?
)
`

type InsertSoldParams struct {
	ProductID   int32  `json:"product_id"`
	ProductName string `json:"product_name"`
	TotalQty    int32  `json:"total_qty"`
	TotalAmount int64  `json:"total_amount"`
}

func (q *Queries) InsertSold(ctx context.Context, arg InsertSoldParams) error {
	_, err := q.db.ExecContext(ctx, insertSold,
		arg.ProductID,
		arg.ProductName,
		arg.TotalQty,
		arg.TotalAmount,
	)
	return err
}
