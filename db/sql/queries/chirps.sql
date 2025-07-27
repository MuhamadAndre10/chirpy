-- name: CreateChirps :one
INSERT INTO chirps (body, user_id)
VALUES ($1, $2)
RETURNING *;
-- name: GetAllChirps :many
SELECT id,
    body,
    user_id,
    created_at,
    updated_at
FROM chirps
ORDER BY created_at ASC;
-- name: GetChirps :one
SELECT id,
    body,
    user_id,
    created_at,
    updated_at
FROM chirps
WHERE id = $1;