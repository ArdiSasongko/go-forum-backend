package contentservice

import (
	"context"
	"database/sql"

	"github.com/ArdiSasongko/go-forum-backend/internal/model"
	"github.com/ArdiSasongko/go-forum-backend/internal/sqlc/comment"
	"github.com/ArdiSasongko/go-forum-backend/internal/sqlc/content"
	"github.com/ArdiSasongko/go-forum-backend/internal/sqlc/useractivities"
	"github.com/sirupsen/logrus"
)

type contentService struct {
	db     *sql.DB
	logger *logrus.Logger
}

func NewContentService(db *sql.DB, logger *logrus.Logger) *contentService {
	return &contentService{
		db:     db,
		logger: logger,
	}
}

type Queries struct {
	ContentQueries        *content.Queries
	CommentQueries        *comment.Queries
	UserActivitiesQueries *useractivities.Queries
}
type ContentService interface {
	InsertContent(ctx context.Context, queries Queries, model model.ContentModel) error
	GetContents(ctx context.Context, queries Queries, limit, offset int32) (*[]model.ContentsResponse, error)
	GetContent(ctx context.Context, queries Queries, contentID, userID, offset, limit int32) (*model.ContentResponse, error)
	UpdateContent(ctx context.Context, queries Queries, contentID, userID int32, req model.UpdateContent) error
	DeleteContent(ctx context.Context, queries Queries, contentID int32) error
	InsertComment(ctx context.Context, queries Queries, req model.CommentModel) error
	DeleteComment(ctx context.Context, queries Queries, userID, contentID int32) error
	LikedDislikeContent(ctx context.Context, queries Queries, req model.LikedModel) error
}
