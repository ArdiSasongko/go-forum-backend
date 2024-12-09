package contentservice

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ArdiSasongko/go-forum-backend/internal/model"
	"github.com/ArdiSasongko/go-forum-backend/internal/sqlc/comment"
	"github.com/ArdiSasongko/go-forum-backend/internal/sqlc/content"
	"github.com/ArdiSasongko/go-forum-backend/utils"
)

func (s *contentService) InsertComment(ctx context.Context, queries Queries, req model.CommentModel) error {
	tx, err := s.db.BeginTx(ctx, nil)
	defer utils.Tx(tx, err)
	commentQueries := queries.CommentQueries.WithTx(tx)
	contentQueries := queries.ContentQueries.WithTx(tx)

	_, err = contentQueries.GetContent(ctx, content.GetContentParams{
		UserID: 0,
		ID:     req.ContentID,
	})
	if err == sql.ErrNoRows {
		s.logger.WithError(err).Error("failed to get content")
		return fmt.Errorf("failed to get content : %v", sql.ErrNoRows.Error())
	} else if err != nil {
		s.logger.WithError(err).Error("failed to get content")
		return fmt.Errorf("failed to get content : %v", err.Error())
	}

	if err := commentQueries.InsertComment(ctx, comment.InsertCommentParams{
		UserID:      req.UserID,
		ContentID:   req.ContentID,
		CommentBody: req.Comment,
		CreatedBy:   req.Username,
		UpdatedBy:   req.Username,
	}); err != nil {
		s.logger.WithError(err).Error("failed to insert comment")
		return fmt.Errorf("failed to insert comment : %v", err.Error())
	}

	s.logger.Info(fmt.Sprintf("success insert comment to contentID : %d", req.ContentID))
	return nil
}

func (s *contentService) DeleteComment(ctx context.Context, queries Queries, userID, contentID int32) error {
	tx, err := s.db.BeginTx(ctx, nil)
	defer utils.Tx(tx, err)
	commentQueries := queries.CommentQueries.WithTx(tx)

	_, err = commentQueries.GetCommentByID(ctx, contentID)
	if err == sql.ErrNoRows {
		s.logger.WithError(err).Error("failed to get content")
		return fmt.Errorf("failed to get content : %v", sql.ErrNoRows.Error())
	} else if err != nil {
		s.logger.WithError(err).Error("failed to get content")
		return fmt.Errorf("failed to get content : %v", err.Error())
	}

	if err := commentQueries.DeleteCommentByUser(ctx, userID); err != nil {
		s.logger.WithError(err).Error("failed to delete comment")
		return fmt.Errorf("failed to delete comment : %v", err.Error())
	}

	s.logger.Info(fmt.Sprintf("success delete comment to contentID : %d", contentID))
	return nil
}
