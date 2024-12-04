package userservice

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ArdiSasongko/go-forum-backend/internal/model"
	"github.com/ArdiSasongko/go-forum-backend/utils"
	"github.com/sirupsen/logrus"
)

func (s *userService) GetProfile(ctx context.Context, username string) (*model.ProfileModel, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	defer utils.Tx(tx, err)

	user, err := s.repo.GetProfile(ctx, tx, username)
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
