package imagesuserrepository

import (
	"context"
	"database/sql"

	imageuser "github.com/ArdiSasongko/go-forum-backend/internal/sqlc/image_user"
)

type imageUserRepository struct {
	queries *imageuser.Queries
}

func NewImageUserRepository(database *sql.DB) *imageUserRepository {
	return &imageUserRepository{queries: imageuser.New(database)}
}

type ImageUserRepository interface {
	GetImage(ctx context.Context, tx *sql.Tx, userID int32) (imageuser.ImagesUser, error)
	InsertImage(ctx context.Context, tx *sql.Tx, model imageuser.CreateImageParams) error
	UpdateImage(ctx context.Context, tx *sql.Tx, model imageuser.UpdateImageParams) error
	DeleteImage(ctx context.Context, tx *sql.Tx, userID int32) error
}

func (r *imageUserRepository) GetImage(ctx context.Context, tx *sql.Tx, userID int32) (imageuser.ImagesUser, error) {
	q := r.queries.WithTx(tx)
	return q.GetImage(ctx, userID)
}

func (r *imageUserRepository) InsertImage(ctx context.Context, tx *sql.Tx, model imageuser.CreateImageParams) error {
	q := r.queries.WithTx(tx)
	return q.CreateImage(ctx, imageuser.CreateImageParams{
		UserID:   model.UserID,
		ImageUrl: model.ImageUrl,
	})
}

func (r *imageUserRepository) UpdateImage(ctx context.Context, tx *sql.Tx, model imageuser.UpdateImageParams) error {
	q := r.queries.WithTx(tx)
	return q.UpdateImage(ctx, imageuser.UpdateImageParams{
		UserID:   model.UserID,
		ImageUrl: model.ImageUrl,
	})
}

func (r *imageUserRepository) DeleteImage(ctx context.Context, tx *sql.Tx, userID int32) error {
	q := r.queries.WithTx(tx)
	return q.DeleteImage(ctx, userID)
}
