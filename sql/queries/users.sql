-- name: CreateUser :one
INSERT INTO users (created_at, updated_at, email)
VALUES ($1, $2, $3)
RETURNING *;