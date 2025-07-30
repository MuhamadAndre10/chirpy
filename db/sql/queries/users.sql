-- name: CreateUser :one
INSERT INTO users (created_at, updated_at, email, hashed_password)
VALUES ($1, $2, $3, $4)
RETURNING *;
-- name: GetUsers :one 
SELECT id,
    email,
    hashed_password,
    created_at,
    updated_at
FROM users
WHERE email = $1;
-- name: CreateRefreshToken :one
INSERT INTO refresh_token (token, expires_at, revoke_at, user_id)
VALUES ($1, $2, $3, $4)
RETURNING *;
-- name: GetRefreshToken :one
SELECT token,
    expires_at,
    revoke_at,
    user_id
FROM refresh_token
WHERE token = $1;