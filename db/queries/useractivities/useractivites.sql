-- name: LikeOrDislikeContent :exec
INSERT INTO user_activities (user_id, content_id, isLiked, isDisliked, created_by, updated_by)
VALUES ($1, $2, $3, $4, $5, $6) 
ON CONFLICT (user_id, content_id) 
DO UPDATE SET
    isLiked = EXCLUDED.isLiked,         
    isDisliked = EXCLUDED.isDisliked,  
    updated_at = current_timestamp,
    updated_by = EXCLUDED.updated_by
WHERE 
    user_activities.isLiked IS DISTINCT FROM EXCLUDED.isLiked 
    OR user_activities.isDisliked IS DISTINCT FROM EXCLUDED.isDisliked;


-- name: GetContentLikes :one
SELECT COUNT(id) FROM user_activities WHERE content_id = $1 AND isLiked = true;

-- name: GetContentDislikes :one
SELECT COUNT(id) FROM user_activities WHERE content_id = $1 AND isDisliked = true;
