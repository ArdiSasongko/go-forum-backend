package contentservice

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

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

	var publicIDs []string
	url := env.GetEnv("CLOUDINARY_URL", "")
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

			var imageURLs []string
			if content.ImageUrls != nil {
				if err := json.Unmarshal(con.ImageUrls, &imageURLs); err != nil {
					errsChan <- fmt.Errorf("failed to parse image urls: %v", err)
					return
				}
			}

			images := make([]model.ImageContent, 0, len(imageURLs))
			for _, url := range imageURLs {
				images = append(images, model.ImageContent{ImageURL: url})
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

	// for _, content := range contents {
	// 	var imageURLs []string
	// 	if content.ImageUrls != nil {
	// 		if err := json.Unmarshal([]byte(content.ImageUrls), &imageURLs); err != nil {
	// 			logrus.WithField("unmarshal image urls", err.Error()).Error("failed to unmarshal image urls")
	// 			return nil, fmt.Errorf("failed to parse image urls: %v", err)
	// 		}
	// 	}
	// 	images := make([]model.ImageContent, 0, len(imageURLs))
	// 	for _, url := range imageURLs {
	// 		images = append(images, model.ImageContent{ImageURL: url})
	// 	}

	// 	item := model.ContentsResponse{
	// 		ContentID:      int(content.ID),
	// 		ContentTitle:   content.ContentTitle,
	// 		ContentBody:    content.ContentBody,
	// 		ContentImage:   images,
	// 		ContentHastags: strings.Split(content.ContentHastags, ","),
	// 	}

	// 	allContents = append(allContents, item)
	// }

	return &allContents, nil
}
