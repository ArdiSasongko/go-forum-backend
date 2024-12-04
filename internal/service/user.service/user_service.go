package userservice

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ArdiSasongko/go-forum-backend/env"
	"github.com/ArdiSasongko/go-forum-backend/internal/model"
	imageuser "github.com/ArdiSasongko/go-forum-backend/internal/sqlc/image_user"
	"github.com/ArdiSasongko/go-forum-backend/internal/sqlc/user"
	cld "github.com/ArdiSasongko/go-forum-backend/pkg/cloudinary"
	"github.com/ArdiSasongko/go-forum-backend/utils"
	"github.com/sirupsen/logrus"
)

func (s *userService) GetProfile(ctx context.Context, email string) (*model.ProfileModel, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	defer utils.Tx(tx, err)

	user, err := s.repo.GetProfile(ctx, tx, email)
	if err == sql.ErrNoRows {
		logrus.WithField("get user", "user didn't exist").Error("user didn't exist")
		return nil, fmt.Errorf("user didn't exist")
	} else if err != nil {
		logrus.WithField("get user", err.Error()).Error(err.Error())
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	response := &model.ProfileModel{
		ID:       user.ID,
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		ImageURL: user.ImageUrl,
		IsValid:  user.IsValid.Bool,
		Role:     string(user.Role),
	}

	return response, nil
}

func (s *userService) UpdateProfile(ctx context.Context, req model.UpdateProfile) error {
	tx, err := s.db.BeginTx(ctx, nil)
	defer utils.Tx(tx, err)

	logrus.Info(req.Email)
	user, err := s.repo.GetUser(ctx, tx, 0, "", req.Email)
	if err == sql.ErrNoRows {
		logrus.WithField("get user", "user didn't exist").Error("user didn't exist")
		return fmt.Errorf("invalid credentials")
	} else if err != nil {
		logrus.WithField("get user", err.Error()).Error(err.Error())
		return fmt.Errorf("failed to get user: %v", err)
	}

	validImga, err := s.imageRepo.GetImage(ctx, tx, user.ID)
	if err != nil {
		logrus.WithField("get image", err.Error()).Error(err.Error())
		return fmt.Errorf("failed to get image: %v", err)
	}

	url := env.GetEnv("CLOUDINARY_URL", "")
	var publicIDs []string
	for _, file := range req.Files {
		imgUrl, publicID, err := cld.UploadImage(ctx, file, url, "forum-profile")
		if err != nil {
			logrus.WithField("upload image", err.Error()).Error(err.Error())
			return fmt.Errorf("failed to upload image profile :%v", err)
		}

		publicIDs = append(publicIDs, publicID)

		imageModel := imageuser.UpdateImageParams{
			UserID:   user.ID,
			ImageUrl: imgUrl,
		}

		err = s.imageRepo.UpdateImage(ctx, tx, imageModel)
		if err != nil {
			for _, id := range publicIDs {
				cld.DestroyImage(ctx, url, id)
			}
			logrus.WithField("update image", err.Error()).Error(err.Error())
			return fmt.Errorf("failed to update image profile :%v", err)
		}
	}

	oldImage, err := cld.GetPublicID(validImga.ImageUrl, "forum-profile")
	if err != nil {
		logrus.WithField("get publicID", err.Error()).Error(err.Error())
		return fmt.Errorf("failed to get publicID: %v", err)
	}

	err = cld.DestroyImage(ctx, url, oldImage)
	if err != nil {
		logrus.WithField("delete image", err.Error()).Error(err.Error())
		return fmt.Errorf("failed to delete image: %v", err)
	}

	return nil
}

func (s *userService) UpdateUser(ctx context.Context, req model.UpdateUser, email string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	defer utils.Tx(tx, err)

	logrus.Info(email)
	validUser, err := s.repo.GetUser(ctx, tx, 0, "", email)
	if err == sql.ErrNoRows {
		logrus.WithField("get user", "user didn't exist").Error("user didn't exist")
		return fmt.Errorf("email didnt exists")
	} else if err != nil {
		logrus.WithField("get user", err.Error()).Error(err.Error())
		return fmt.Errorf("failed to get user: %v", err)
	}

	modelUpdate := user.UpdateUserParams{
		Name:     utils.DefaultValue[string](validUser.Name, req.Name),
		Username: utils.DefaultValue[string](validUser.Username, req.Username),
		ID:       validUser.ID,
	}

	if err := s.repo.UpdateUser(ctx, tx, modelUpdate); err != nil {
		logrus.WithField("update user", err.Error()).Error(err.Error())
		return fmt.Errorf("failed to update user : %v", err)
	}

	return nil
}
