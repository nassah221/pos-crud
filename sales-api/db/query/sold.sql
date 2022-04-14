-- name: InsertSold :exec
INSERT INTO sold (
    product_id, product_name, total_qty, total_amount
) VALUES (
    ?,?,?,?
);

-- name: GetSold :many
SELECT id,product_id,product_name,SUM(total_qty) AS total_qty,SUM(total_amount) AS total_Amount 
FROM sold GROUP BY product_id ORDER BY id;

-- name: SetSQLMode :execresult
SET SESSION sql_mode=(SELECT REPLACE(@@sql_mode,'ONLY_FULL_GROUP_BY','')); SELECT @@sql_mode;