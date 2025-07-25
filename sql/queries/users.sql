-- name: CreateUser :one
INSERT INTO users (created_at, updated_at, email)
VALUES (now(), now(), "andrepriyanto95@gmail.com")
RETURNING *;