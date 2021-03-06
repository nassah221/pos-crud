-- name: CreateDiscount :execresult
INSERT INTO discounts (
    qty, type, result, expired_at, expired_at_format, string_format
) VALUES (
    ?,?,?,?,?,?
);

-- name: GetDiscount :one
SELECT * FROM discounts 
WHERE id=? LIMIT 1;

-- name: CreateProductDiscount :exec
INSERT INTO product_discount (
    discount_id, product_id
) VALUES (
    ?,?
);