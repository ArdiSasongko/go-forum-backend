package contentservice

import (
	"context"
	"database/sql"

	"github.com/ArdiSasongko/go-forum-backend/internal/model"
	"github.com/ArdiSasongko/go-forum-backend/internal/sqlc/content"
)

type contentService struct {
	db *sql.DB
}

func NewContentService(db *sql.DB) *contentService {
	return &contentService{
		db: db,
	}
}

type Queries struct {
	ContentQueries *content.Queries
}
type ContentService interface {
	InsertContent(ctx context.Context, queries Queries, model model.ContentModel) error
	GetContents(ctx context.Context, queries Queries, limit, offset int32) (*[]model.ContentsResponse, error)
}
