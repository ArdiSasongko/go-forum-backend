-- name: GetUser :one
SELECT * FROM users WHERE id = $1 OR username = $2 OR email = $3;

-- name: CreateUser :one
INSERT INTO users (name, username, email, password, role, is_valid)
VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;

-- name: UpdateUser :exec
UPDATE users set name = $1, username = $2 where id = $3;

-- name: UpdatePassword :exec
UPDATE users set password = $1 where id = $2;

-- name: ValidateUser :exec
UPDATE users set is_valid = true where id = $1;

-- name: GetUserProfile :one
SELECT u.id, u.name, u.username, u.email, i.image_url, u.is_valid, u.role FROM users u JOIN images_user i ON u.id = i.user_id WHERE u.email = $1;