-- name: InsertComment :exec
INSERT INTO comments (user_id, content_id, comment_body, created_by, updated_by) VALUES ($1, $2, $3, $4, $5);

-- name: GetCommentByContent :many
SELECT * FROM comments WHERE content_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: GetCommentByUser :many
SELECT * FROM comments WHERE user_id = $1;

-- name: GetCommentByID :one
SELECT * FROM comments WHERE id = $1;

-- name: DeleteCommentByUser :exec
DELETE FROM comments WHERE user_id = $1;

-- name: GetCountOfComments :one
SELECT COUNT(id) FROM comments WHERE content_id = $1;