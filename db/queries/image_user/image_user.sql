-- name: GetImage :one
SELECT * FROM images_user WHERE user_id = $1;

-- name: CreateImage :exec
INSERT INTO images_user (user_id, image_url) VALUES ($1, $2);

-- name: DeleteImage :exec
DELETE FROM images_user WHERE user_id = $1;

-- name: UpdateImage :exec
UPDATE images_user SET image_url = $2 WHERE user_id = $1;