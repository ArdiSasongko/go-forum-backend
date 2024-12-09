package userservice

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ArdiSasongko/go-forum-backend/internal/model"
	tokentable "github.com/ArdiSasongko/go-forum-backend/internal/sqlc/token"
	"github.com/ArdiSasongko/go-forum-backend/internal/sqlc/user"
	"github.com/ArdiSasongko/go-forum-backend/utils"
	"github.com/sirupsen/logrus"
)

func (s *userService) ValidateEmail(ctx context.Context, queries Queries, payload model.ValidatePayload) error {
	tx, err := s.db.BeginTx(ctx, nil)
	defer utils.Tx(tx, err)

	userQueries := queries.UserQueries.WithTx(tx)
	tokenQueries := queries.TokenQueries.WithTx(tx)

	logrus.Info("username", payload.Username)
	user, err := userQueries.GetUser(ctx, user.GetUserParams{
		ID:       0,
		Username: payload.Username,
		Email:    "",
	})
	if err == sql.ErrNoRows {
		s.logger.WithError(err).Error("user didnt exists")
		return fmt.Errorf("user didn't exist")
	} else if err != nil {
		s.logger.WithError(err).Error("failed to get user")
		return fmt.Errorf("failed to get user: %v", err)
	}

	validToken, err := tokenQueries.GetToken(ctx, tokentable.GetTokenParams{
		UserID: user.ID,
		Token:  payload.Token,
	})
	if err == sql.ErrNoRows {
		s.logger.WithError(err).Error("user didnt exists")
		return fmt.Errorf("user didn't exist")
	} else if err != nil {
		s.logger.WithError(err).Error("failed to get user")
		return fmt.Errorf("failed to get user: %v", err)
	}

	if validToken.ExpiredAt.Before(time.Now().UTC()) {
		s.logger.WithError(err).Error("token has expired")
		return fmt.Errorf("token has expired, please resend new token")
	}

	if err := userQueries.ValidateUser(ctx, user.ID); err != nil {
		s.logger.WithError(err).Error("failed to validate user")
		return fmt.Errorf("failed to validate user : %v", err)
	}

	if err := tokenQueries.DeleteToken(ctx, tokentable.DeleteTokenParams{
		UserID:    user.ID,
		TokenType: "email",
	}); err != nil {
		s.logger.WithError(err).Error("failed to delete token")
		return fmt.Errorf("failed to delete token : %v", err)
	}

	s.logger.Info(fmt.Sprintf("user %v validation email", user.ID))
	return nil
}

func (s *userService) ResendEmail(ctx context.Context, queries Queries, payload model.ValidatePayload) error {
	tx, err := s.db.BeginTx(ctx, nil)
	defer utils.Tx(tx, err)

	userQueries := queries.UserQueries.WithTx(tx)
	tokenQueries := queries.TokenQueries.WithTx(tx)

	user, err := userQueries.GetUser(ctx, user.GetUserParams{
		ID:       0,
		Username: payload.Username,
		Email:    "",
	})
	if err == sql.ErrNoRows {
		s.logger.WithError(err).Error("user didnt exists")
		return fmt.Errorf("user didn't exist")
	} else if err != nil {
		s.logger.WithError(err).Error("failed to get user")
		return fmt.Errorf("failed to get user: %v", err)
	}

	validationToken := utils.GenToken()
	err = tokenQueries.UpdateToken(ctx, tokentable.UpdateTokenParams{
		UserID:    user.ID,
		Token:     int32(validationToken),
		ExpiredAt: time.Now().UTC().Add(5 * time.Minute),
	})
	if err != nil {
		s.logger.WithError(err).Error("failed to update token")
		return fmt.Errorf("failed to create new token :%v", err)
	}

	logrus.Info(validationToken)
	err = utils.SendToken(user.Email, "resend_email", int32(validationToken))
	if err != nil {
		s.logger.WithError(err).Error("failed to resend email")
		return fmt.Errorf("failed to send email :%v", err)
	}

	s.logger.Info(fmt.Sprintf("user %v success get new token", user.ID))
	return nil
}
