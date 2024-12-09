package contentservice

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ArdiSasongko/go-forum-backend/internal/model"
	"github.com/ArdiSasongko/go-forum-backend/internal/sqlc/content"
	"github.com/ArdiSasongko/go-forum-backend/internal/sqlc/useractivities"
	"github.com/ArdiSasongko/go-forum-backend/utils"
	"github.com/sirupsen/logrus"
)

func (s *contentService) LikedDislikeContent(ctx context.Context, queries Queries, req model.LikedModel) error {
	tx, err := s.db.BeginTx(ctx, nil)
	defer utils.Tx(tx, err)

	logrus.Info(req.IsLike, req.IsDislike)
	userActivitiesQueries := queries.UserActivitiesQueries.WithTx(tx)
	contentQueries := queries.ContentQueries.WithTx(tx)

	if req.IsDislike && req.IsLike {
		s.logger.Error("cant do two activities")
		return fmt.Errorf("cant do two activities")
	}
	_, err = contentQueries.GetContent(ctx, content.GetContentParams{
		UserID: 0,
		ID:     req.ContentID,
	})
	if err == sql.ErrNoRows {
		s.logger.WithError(err).Error("content didnt exist")
		return fmt.Errorf("failed to get content: %v", err)
	} else if err != nil {
		s.logger.WithError(err).Error("failed to get content")
		return fmt.Errorf("failed to get content: %v", err)
	}

	if err := userActivitiesQueries.LikeOrDislikeContent(ctx, useractivities.LikeOrDislikeContentParams{
		UserID:    req.UserID,
		ContentID: req.ContentID,
		CreatedBy: req.Username,
		UpdatedBy: req.Username,
		Isdisliked: sql.NullBool{
			Bool:  req.IsDislike,
			Valid: true,
		},
		Isliked: sql.NullBool{
			Bool:  req.IsLike,
			Valid: true,
		},
	}); err != nil {
		s.logger.WithError(err).Error("failed to like content")
		return fmt.Errorf("failed to like content : %v", err)
	}

	s.logger.Info(fmt.Sprintf("user %d success liked content %d", req.UserID, req.ContentID))
	return nil
}
