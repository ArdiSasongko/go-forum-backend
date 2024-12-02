-- name: GetToken :one
SELECT * FROM user_sessions WHERE user_id = $1;

-- name: InsertToken :one
INSERT INTO user_sessions (user_id, token, token_expired, refresh_token, refresh_token_expired) VALUES ($1, $2, $3, $4, $5) RETURNING token, refresh_token;

-- name: UpdateToken :exec
UPDATE user_sessions set token = $1, token_expired = $2, refresh_token = $3, refresh_token_expired = $4 where user_id = $5;

-- name: DeleteToken :exec
DELETE FROM user_sessions WHERE user_id = $1;