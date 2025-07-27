-- name: CreateUser :one
INSERT INTO users (created_at, updated_at, email)
VALUES ($1, $2, $3)
RETURNING *;
-- name: CreateChirps :one
INSERT INTO chirps (body, user_id)
VALUES ($1, $2)
RETURNING *;