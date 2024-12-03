package userservice

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ArdiSasongko/go-forum-backend/internal/model"
	tokentable "github.com/ArdiSasongko/go-forum-backend/internal/sqlc/token"
	"github.com/ArdiSasongko/go-forum-backend/utils"
	"github.com/sirupsen/logrus"
)

func (s *userService) ValidateEmail(ctx context.Context, payload model.ValidatePayload) error {
	tx, err := s.db.BeginTx(ctx, nil)
	defer utils.Tx(tx, err)

	logrus.Info("username", payload.Username)
	user, err := s.repo.GetUser(ctx, tx, 0, payload.Username, "")
	if err == sql.ErrNoRows {
		logrus.WithField("get user", "user didn't exist").Error("user didn't exist")
		return fmt.Errorf("invalid credentials")
	} else if err != nil {
		logrus.WithField("get user", err.Error()).Error(err.Error())
		return fmt.Errorf("failed to get user: %v", err)
	}

	validToken, err := s.tokenRepo.GetToken(ctx, tx, user.ID, payload.Token)
	if err != nil {
		logrus.WithField("get validate token", err.Error()).Error(err.Error())
		return fmt.Errorf("failed get validate token : %v", err)
	} else if err == sql.ErrNoRows {
		logrus.WithField("get validate token", "token didnt exists please resend again").Error("token didnt exists please resend again")
		return fmt.Errorf("token didnt exists please resend again")
	}

	if validToken.ExpiredAt.Before(time.Now().UTC()) {
		logrus.WithField("get validate token", "token has expired").Error("token has expired")
		return fmt.Errorf("token has expired, please resend new token")
	}

	if err := s.repo.ValidateUser(ctx, tx, user.ID); err != nil {
		logrus.WithField("validate user", err.Error()).Error(err.Error())
		return fmt.Errorf("failed to validate user : %v", err)
	}

	if err := s.tokenRepo.DeleteToken(ctx, tx, user.ID, "email"); err != nil {
		logrus.WithField("delete token", err.Error()).Error(err.Error())
		return fmt.Errorf("failed to delete token : %v", err)
	}

	return nil
}

func (s *userService) ResendEmail(ctx context.Context, payload model.ValidatePayload) error {
	tx, err := s.db.BeginTx(ctx, nil)
	defer utils.Tx(tx, err)

	user, err := s.repo.GetUser(ctx, tx, 0, payload.Username, "")
	if err == sql.ErrNoRows {
		logrus.WithField("get user", "user didn't exist").Error("user didn't exist")
		return fmt.Errorf("invalid credentials")
	} else if err != nil {
		logrus.WithField("get user", err.Error()).Error(err.Error())
		return fmt.Errorf("failed to get user: %v", err)
	}

	validationToken := utils.GenToken()
	tokenModel := tokentable.UpdateTokenParams{
		UserID:    user.ID,
		Token:     int32(validationToken),
		ExpiredAt: time.Now().UTC().Add(5 * time.Minute),
	}

	err = s.tokenRepo.UpdateToken(ctx, tx, tokenModel)
	if err != nil {
		logrus.WithField("update token", err.Error()).Error(err.Error())
		return fmt.Errorf("failed to create new token :%v", err)
	}

	logrus.Info(validationToken)
	err = utils.SendToken(user.Email, "resend_email", int32(validationToken))
	if err != nil {
		logrus.WithField("send email", err.Error()).Error(err.Error())
		return fmt.Errorf("failed to send email :%v", err)
	}

	return nil
}
