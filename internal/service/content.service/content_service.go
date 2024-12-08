package contentservice

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ArdiSasongko/go-forum-backend/env"
	"github.com/ArdiSasongko/go-forum-backend/internal/model"
	"github.com/ArdiSasongko/go-forum-backend/internal/sqlc/content"
	cld "github.com/ArdiSasongko/go-forum-backend/pkg/cloudinary"
	"github.com/ArdiSasongko/go-forum-backend/utils"
	"github.com/sirupsen/logrus"
)

func (s *contentService) InsertContent(ctx context.Context, queries Queries, model model.ContentModel) error {
	tx, err := s.db.BeginTx(ctx, nil)
	defer utils.Tx(tx, err)

	contentQueries := queries.ContentQueries.WithTx(tx)

	contentID, err := contentQueries.InsertContent(ctx, content.InsertContentParams{
		UserID:         model.UserID,
		ContentTitle:   model.ContentTitle,
		ContentBody:    model.ContentBody,
		ContentHastags: strings.Join(model.ContentHastags, ","),
		CreatedBy:      model.Username,
		UpdatedBy:      model.Username,
	})
	if err != nil {
		logrus.WithField("insert content", err.Error()).Error(err.Error())
		return fmt.Errorf("failed to to insert content : %v", err)
	}

	var (
		publicIDs []string
		url       = env.GetEnv("CLOUDINARY_URL", "")
	)

	for _, image := range model.Files {
		imgUrl, publicID, err := cld.UploadImage(ctx, image, url, "forum-content")
		if err != nil {
			logrus.WithField("upload content", err.Error()).Error(err.Error())
			return fmt.Errorf("failed to to upload content : %v", err)
		}

		publicIDs = append(publicIDs, publicID)

		if err := contentQueries.InsertImageContent(ctx, content.InsertImageContentParams{
			ContentID: contentID,
			ImageUrl:  imgUrl,
		}); err != nil {
			for _, id := range publicIDs {
				cld.DestroyImage(ctx, url, id)
			}

			if err := contentQueries.DeleteImageContent(ctx, content.DeleteImageContentParams{
				ContentID: contentID,
				ImageUrl:  imgUrl,
			}); err != nil {
				logrus.WithField("delete content image", err.Error()).Error(err.Error())
				return fmt.Errorf("failed to delete content image : %v", err)
			}
			logrus.WithField("insert content image", err.Error()).Error(err.Error())
			return fmt.Errorf("failed to insert content image : %v", err)
		}
	}
	return nil
}

func (s *contentService) GetContents(ctx context.Context, queries Queries, limit, offset int32) (*[]model.ContentsResponse, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	defer utils.Tx(tx, err)

	contentQueries := queries.ContentQueries.WithTx(tx)

	contents, err := contentQueries.GetContents(ctx, content.GetContentsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		logrus.WithField("get contents", err.Error()).Error(err.Error())
		return nil, fmt.Errorf("failed to get contents : %v", err)
	}

	if len(contents) == 0 {
		logrus.WithField("get contents", "contents didnt exists").Error("contents didnt exists")
		return nil, fmt.Errorf("contents didnt exists")
	}

	allContents := make([]model.ContentsResponse, 0, len(contents))
	var wg sync.WaitGroup
	resultsChan := make(chan model.ContentsResponse, len(contents))
	errsChan := make(chan error, len(contents))

	for _, con := range contents {
		wg.Add(1)
		go func(content content.GetContentsRow) {
			defer wg.Done()

			var images []model.ImageContent
			if len(content.ImageUrls) > 0 {
				imageURLs := strings.Split(string(content.ImageUrls), ",")
				for _, url := range imageURLs {
					images = append(images, model.ImageContent{ImageURL: strings.TrimSpace(url)})
				}
			}

			item := model.ContentsResponse{
				ContentID:      int(content.ID),
				ContentTitle:   content.ContentTitle,
				ContentBody:    content.ContentBody,
				ContentImage:   images,
				ContentHastags: strings.Split(content.ContentHastags, ","),
			}

			resultsChan <- item
		}(con)
	}

	go func() {
		wg.Wait()
		close(resultsChan)
		close(errsChan)
	}()

	for err := range errsChan {
		logrus.WithField("error", err.Error()).Error("error process content")
		return nil, err
	}

	for item := range resultsChan {
		allContents = append(allContents, item)
	}

	return &allContents, nil
}

func (s *contentService) GetContent(ctx context.Context, queries Queries, contentID int32) (*model.ContentResponse, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	defer utils.Tx(tx, err)
	contentQueries := queries.ContentQueries.WithTx(tx)

	contentRow, err := contentQueries.GetContent(ctx, contentID)
	if err != nil {
		if err == sql.ErrNoRows {
			logrus.WithField("get content", "content didnt exist").Error("failed to get content")
			return nil, fmt.Errorf("failed to get content : %v", err)
		}
		logrus.WithField("get content", err.Error()).Error("failed to get content")
		return nil, fmt.Errorf("failed to get content : %v", err)
	}

	contentResponse := model.ContentResponse{
		ContentID:      int(contentRow.ID),
		ContentTitle:   contentRow.ContentTitle,
		ContentBody:    contentRow.ContentBody,
		ContentHastags: strings.Split(contentRow.ContentHastags, ","),
		CreatedAt:      contentRow.CreatedAt,
		UpdatedAt:      contentRow.UpdatedAt,
		CreatedBy:      contentRow.CreatedBy,
	}

	logrus.Info(contentResponse.ContentID)
	allImages := make([]model.ImageContent, 0, len(contentRow.ImageUrls))
	var wg sync.WaitGroup
	resultsChan := make(chan model.ImageContent, len(contentRow.ImageUrls))

	if len(contentRow.ImageUrls) > 0 {
		imageURLs := strings.Split(string(contentRow.ImageUrls), ",")
		for _, url := range imageURLs {
			wg.Add(1)
			go func(imageURL string) {
				defer wg.Done()
				resultsChan <- model.ImageContent{ImageURL: strings.TrimSpace(imageURL)}
			}(url)
		}
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	go func() {
		for img := range resultsChan {
			allImages = append(allImages, img)
		}
	}()

	wg.Wait()

	contentResponse.ContentImage = allImages

	return &contentResponse, nil
}

func (s *contentService) UpdateContent(ctx context.Context, queries Queries, contentID, userID int32, req model.UpdateContent) error {
	tx, err := s.db.BeginTx(ctx, nil)
	defer utils.Tx(tx, err)
	contentQueries := queries.ContentQueries.WithTx(tx)

	oldContent, err := s.GetContent(ctx, queries, contentID)
	if err != nil {
		logrus.WithField("get content", err.Error()).Error("failed to get content")
		return fmt.Errorf("failed to get content : %v", err)
	}

	newHastags := strings.Join(req.ContentHastags, ",")
	oldHastags := strings.Join(oldContent.ContentHastags, ",")
	if err := contentQueries.UpdateContent(ctx, content.UpdateContentParams{
		ID:             contentID,
		ContentTitle:   utils.DefaultValue[string](oldContent.ContentTitle, req.ContentTitle),
		ContentBody:    utils.DefaultValue[string](oldContent.ContentBody, req.ContentBody),
		ContentHastags: utils.DefaultValue[string](oldHastags, newHastags),
		UpdatedBy:      req.UpdatedBy,
		UpdatedAt:      time.Now().UTC(),
		UserID:         userID,
	}); err != nil {
		logrus.WithField("update content", err.Error()).Error("failed to update content")
		return fmt.Errorf("failed to update content : %v", err)
	}

	return nil
}

func (s *contentService) DeleteContent(ctx context.Context, queries Queries, contentID int32) error {
	tx, err := s.db.BeginTx(ctx, nil)
	defer utils.Tx(tx, err)
	contentQueries := queries.ContentQueries.WithTx(tx)

	validContent, err := s.GetContent(ctx, queries, contentID)
	if err != nil {
		logrus.WithField("get content", err.Error()).Error("failed to get content")
		return fmt.Errorf("failed to get content : %v", err)
	}

	if err := contentQueries.DeleteContent(ctx, contentID); err != nil {
		logrus.WithField("delete content", err.Error()).Error("failed to delete content")
		return fmt.Errorf("failed to delete content : %v", err)
	}

	var (
		publicIDs []string
		url       = env.GetEnv("CLOUDINARY_URL", "")
	)

	for _, url := range validContent.ContentImage {
		publicID, _ := cld.GetPublicID(url.ImageURL, "forum-content")
		publicIDs = append(publicIDs, publicID)
	}

	for _, publicId := range publicIDs {
		logrus.Info(publicId)
		if err := cld.DestroyImage(ctx, url, publicId); err != nil {
			logrus.WithField("delete image", err.Error()).Error("failed to delete image content")
			return fmt.Errorf("failed to delete image content : %v", err)
		}
	}

	return nil
}
