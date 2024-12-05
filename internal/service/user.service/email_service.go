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
		logrus.WithField("get user", "user didn't exist").Error("user didn't exist")
		return fmt.Errorf("invalid credentials")
	} else if err != nil {
		logrus.WithField("get user", err.Error()).Error(err.Error())
		return fmt.Errorf("failed to get user: %v", err)
	}

	validToken, err := tokenQueries.GetToken(ctx, tokentable.GetTokenParams{
		UserID: user.ID,
		Token:  payload.Token,
	})
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

	if err := userQueries.ValidateUser(ctx, user.ID); err != nil {
		logrus.WithField("validate user", err.Error()).Error(err.Error())
		return fmt.Errorf("failed to validate user : %v", err)
	}

	if err := tokenQueries.DeleteToken(ctx, tokentable.DeleteTokenParams{
		UserID:    user.ID,
		TokenType: "email",
	}); err != nil {
		logrus.WithField("delete token", err.Error()).Error(err.Error())
		return fmt.Errorf("failed to delete token : %v", err)
	}

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
		logrus.WithField("get user", "user didn't exist").Error("user didn't exist")
		return fmt.Errorf("invalid credentials")
	} else if err != nil {
		logrus.WithField("get user", err.Error()).Error(err.Error())
		return fmt.Errorf("failed to get user: %v", err)
	}

	validationToken := utils.GenToken()
	err = tokenQueries.UpdateToken(ctx, tokentable.UpdateTokenParams{
		UserID:    user.ID,
		Token:     int32(validationToken),
		ExpiredAt: time.Now().UTC().Add(5 * time.Minute),
	})
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
