// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: comment.sql

package comment

import (
	"context"
)

const deleteCommentByUser = `-- name: DeleteCommentByUser :exec
DELETE FROM comments WHERE user_id = $1
`

func (q *Queries) DeleteCommentByUser(ctx context.Context, userID int32) error {
	_, err := q.db.ExecContext(ctx, deleteCommentByUser, userID)
	return err
}

const getCommentByContent = `-- name: GetCommentByContent :many
SELECT id, user_id, content_id, comment_body, created_at, updated_at, created_by, updated_by FROM comments WHERE content_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3
`

type GetCommentByContentParams struct {
	ContentID int32
	Limit     int32
	Offset    int32
}

func (q *Queries) GetCommentByContent(ctx context.Context, arg GetCommentByContentParams) ([]Comment, error) {
	rows, err := q.db.QueryContext(ctx, getCommentByContent, arg.ContentID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Comment
	for rows.Next() {
		var i Comment
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.ContentID,
			&i.CommentBody,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.CreatedBy,
			&i.UpdatedBy,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getCommentByID = `-- name: GetCommentByID :one
SELECT id, user_id, content_id, comment_body, created_at, updated_at, created_by, updated_by FROM comments WHERE id = $1
`

func (q *Queries) GetCommentByID(ctx context.Context, id int32) (Comment, error) {
	row := q.db.QueryRowContext(ctx, getCommentByID, id)
	var i Comment
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ContentID,
		&i.CommentBody,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.CreatedBy,
		&i.UpdatedBy,
	)
	return i, err
}

const getCommentByUser = `-- name: GetCommentByUser :many
SELECT id, user_id, content_id, comment_body, created_at, updated_at, created_by, updated_by FROM comments WHERE user_id = $1
`

func (q *Queries) GetCommentByUser(ctx context.Context, userID int32) ([]Comment, error) {
	rows, err := q.db.QueryContext(ctx, getCommentByUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Comment
	for rows.Next() {
		var i Comment
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.ContentID,
			&i.CommentBody,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.CreatedBy,
			&i.UpdatedBy,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getCountOfComments = `-- name: GetCountOfComments :one
SELECT COUNT(id) FROM comments WHERE content_id = $1
`

func (q *Queries) GetCountOfComments(ctx context.Context, contentID int32) (int64, error) {
	row := q.db.QueryRowContext(ctx, getCountOfComments, contentID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const insertComment = `-- name: InsertComment :exec
INSERT INTO comments (user_id, content_id, comment_body, created_by, updated_by) VALUES ($1, $2, $3, $4, $5)
`

type InsertCommentParams struct {
	UserID      int32
	ContentID   int32
	CommentBody string
	CreatedBy   string
	UpdatedBy   string
}

func (q *Queries) InsertComment(ctx context.Context, arg InsertCommentParams) error {
	_, err := q.db.ExecContext(ctx, insertComment,
		arg.UserID,
		arg.ContentID,
		arg.CommentBody,
		arg.CreatedBy,
		arg.UpdatedBy,
	)
	return err
}
