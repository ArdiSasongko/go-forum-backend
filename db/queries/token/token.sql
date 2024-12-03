-- name: GetToken :one
SELECT * FROM tokens WHERE user_id = $1 AND token = $2;

-- name: CreateToken :exec
INSERT INTO tokens (user_id, token_type, token, expired_at) VALUES ($1, $2, $3, $4);

-- name: UpdateToken :exec
UPDATE tokens set token = $1, expired_at = $2 where user_id = $3;

-- name: DeleteToken :exec
Delete FROM tokens WHERE user_id = $1 AND token_type = $2;
