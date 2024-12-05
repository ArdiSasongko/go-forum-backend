-- name: InsertContent :one
INSERT INTO 
contents (user_id, content_title, content_body, content_hastags, created_by, updated_by) 
VALUES 
($1, $2, $3, $4, $5, $6) RETURNING id;

-- name: GetContent :many
SELECT 
c.id, 
c.content_title,
c.content_body,
ARRAY_AGG(i.image_url) as image_urls, 
c.content_hastags, 
c.created_at, 
c.updated_at, 
c.created_by 
FROM contents c 
LEFT JOIN images_content i 
ON c.id = i.content_id 
WHERE c.id = $1;

-- name: GetContents :many
SELECT 
c.id, 
c.content_title, 
c.content_body,
c.content_hastags, 
(
    SELECT json_agg(image_url)
    FROM images_content
    WHERE content_id = c.id
) AS image_urls
FROM contents c 
LEFT JOIN images_content i 
ON c.id = i.content_id 
GROUP BY c.id
ORDER BY c.created_at 
DESC LIMIT $1 OFFSET $2;

-- name: UpdateContent :exec
UPDATE contents SET content_title = $2, content_body = $3, content_hastags = $4, updated_by = $5 WHERE id = $1 AND user_id = $6;

-- name: DeleteContent :exec
DELETE FROM contents WHERE id = $1;

-- name: InsertImageContent :exec
INSERT INTO images_content (content_id, image_url) VALUES ($1, $2);

-- name: GetImagesContent :many
SELECT image_url FROM images_content WHERE content_id = $1;

-- name: DeleteImagesContent :exec
DELETE FROM images_content WHERE content_id = $1;

-- name: DeleteImageContent :exec
DELETE FROM images_content WHERE content_id = $1 AND image_url = $2;

-- name: UpdateImageContent :exec
UPDATE images_content SET image_url = $2 WHERE content_id = $1;