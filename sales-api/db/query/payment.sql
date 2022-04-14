-- name: ListPayments :many
SELECT * FROM payments
ORDER BY id LIMIT ? OFFSET ?;

-- name: GetPayment :one
SELECT * FROM payments
WHERE id = ? LIMIT 1;

-- name: CreatePayment :execresult
INSERT INTO payments (
    name, type, logo
) VALUES (
    ?,?,?
);

-- name: UpdatePayment :exec
UPDATE payments SET name = ?,type=?,logo=?,updated_at=CURRENT_TIMESTAMP
WHERE id = ?;

-- name: DeletePayment :exec
DELETE FROM payments
WHERE id = ?;