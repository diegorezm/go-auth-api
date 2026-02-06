-- name: CreateUser :one
INSERT INTO users (
    email,
    password_hash,
    name,
    avatar_url
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1
LIMIT 1;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1
LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET
    email = $2,
    name = $3,
    avatar_url = $4,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteUserById :exec
DELETE FROM users
WHERE id = $1;
