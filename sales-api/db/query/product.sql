-- name: CreateProduct :execresult
INSERT INTO products (
    name,image_url,price,stock,category_id,sku
) VALUES (
    ?,?,?,?,?,?
);

-- name: ListProducts :many
SELECT * FROM products
ORDER BY id LIMIT ? OFFSET ?;

-- name: GetProduct :one
SELECT * FROM products
WHERE id = ? LIMIT 1;

-- name: UpdateProduct :exec
UPDATE products SET category_id=?,name=?,image_url=?,price=?,stock=?,updated_at=CURRENT_TIMESTAMP
WHERE id = ?;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = ?;

-- name: ListProductsDetailsByCategoryID :many
select p.id,p.name,p.sku,p.image_url,p.stock,p.price,p.category_id,
c.name as category_name,
dis.id as discount_id,
d.qty,d.type,d.result,d.expired_at,d.expired_at_format,d.string_format
from products as p left join categories as c on c.id = p.category_id
left join product_discount as dis on dis.product_id = p.id
left join discounts as d on d.id=discount_id
where category_id=? ORDER BY p.id LIMIT ? OFFSET ?;

-- name: ListProductsDetailsByName :many
select p.id,p.name,p.sku,p.image_url,p.stock,p.price,p.category_id,
c.name as category_name,
dis.id as discount_id,
d.qty,d.type,d.result,d.expired_at,d.expired_at_format,d.string_format
from products as p left join categories as c on c.id = p.category_id
left join product_discount as dis on dis.product_id = p.id
left join discounts as d on d.id=discount_id
WHERE p.name LIKE ? ORDER BY p.id LIMIT ? OFFSET ?;

-- name: GetProductDetails :one
SELECT p.name,p.id,p.sku,p.stock,p.price,p.image_url,p.category_id,
dis.id AS discount_id,dis.qty,dis.type,dis.expired_at,dis.result,dis.expired_at_format,dis.string_format,
c.name AS category_name FROM products AS p JOIN categories AS c ON c.id = p.category_id 
LEFT JOIN product_discount AS d ON p.id = d.product_id LEFT JOIN discounts AS dis ON dis.id=discount_id 
WHERE p.id = ?;