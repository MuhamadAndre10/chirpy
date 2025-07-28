-- name: CreateUser :one
INSERT INTO users (created_at, updated_at, email, hashed_password)
VALUES ($1, $2, $3, $4)
RETURNING *;
-- name: GetUsers :one 
SELECT email,
    hashed_password,
    created_at,
    updated_at
FROM users
WHERE email = $1;