-- name: ListCategories :many
SELECT * FROM categories
ORDER BY id LIMIT ? OFFSET ?;

-- name: GetCategory :one
SELECT * FROM categories
WHERE id = ? LIMIT 1;

-- name: CreateCategory :execresult
INSERT INTO categories (
    name
) VALUES (
    ?
);

-- name: UpdateCategory :exec
UPDATE categories SET name = ?,updated_at=CURRENT_TIMESTAMP
WHERE id = ?;

-- name: DeleteCategory :exec
DELETE FROM categories
WHERE id = ?;