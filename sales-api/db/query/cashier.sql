-- name: CreateCashier :execresult
INSERT INTO cashiers (
  name, password
) VALUES (
  ?, ?
);

-- name: GetCashier :one
SELECT cashiers.password, cashiers.name FROM cashiers
WHERE id = ? LIMIT 1;

-- name: ListCashiers :many
SELECT * FROM cashiers
ORDER BY id LIMIT ? OFFSET ?;

-- name: UpdateCashier :exec
UPDATE cashiers SET name=?,password=?,updated_at=CURRENT_TIMESTAMP 
WHERE id = ?;

-- name: DeleteCashier :exec
DELETE FROM cashiers
WHERE id = ?;

-- name: DetailCashier :one
SELECT * FROM cashiers
WHERE id = ? LIMIT 1;
