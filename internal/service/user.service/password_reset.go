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
	"golang.org/x/crypto/bcrypt"
)

func (s *userService) ResetPassword(ctx context.Context, queries Queries, req model.SendEmail) error {
	tx, err := s.db.BeginTx(ctx, nil)
	defer utils.Tx(tx, err)
	userQueries := queries.UserQueries.WithTx(tx)
	tokenQueries := queries.TokenQueries.WithTx(tx)

	user, err := userQueries.GetUser(ctx, user.GetUserParams{
		ID:       0,
		Username: "",
		Email:    req.Email,
	})
	if err == sql.ErrNoRows {
		s.logger.WithError(err).Error("user didnt exists")
		return fmt.Errorf("user didn't exist")
	} else if err != nil {
		s.logger.WithError(err).Error("failed to get user")
		return fmt.Errorf("failed to get user: %v", err)
	}

	validationToken := utils.GenToken()
	err = tokenQueries.CreateToken(ctx, tokentable.CreateTokenParams{
		UserID:    user.ID,
		TokenType: "password_reset",
		Token:     int32(validationToken),
		ExpiredAt: time.Now().UTC().Add(5 * time.Minute),
	})
	if err != nil {
		s.logger.WithError(err).Error("failed to create token")
		return fmt.Errorf("failed to create new token :%v", err)
	}

	logrus.Info(validationToken)
	err = utils.SendToken(req.Email, "password", int32(validationToken))
	if err != nil {
		s.logger.WithError(err).Error("failed to send email")
		return fmt.Errorf("failed to send email :%v", err)
	}

	s.logger.Info(fmt.Sprintf("user %v get reset token", user.ID))
	return nil
}

func (s *userService) ConfirmPassword(ctx context.Context, queries Queries, req model.ResetPassword) error {
	tx, err := s.db.BeginTx(ctx, nil)
	defer utils.Tx(tx, err)
	userQueries := queries.UserQueries.WithTx(tx)
	tokenQueries := queries.TokenQueries.WithTx(tx)

	validUser, err := userQueries.GetUser(ctx, user.GetUserParams{
		ID:       0,
		Username: "",
		Email:    req.Email,
	})
	if err == sql.ErrNoRows {
		s.logger.WithError(err).Error("user didnt exists")
		return fmt.Errorf("user didn't exist")
	} else if err != nil {
		s.logger.WithError(err).Error("failed to get user")
		return fmt.Errorf("failed to get user: %v", err)
	}

	validToken, err := tokenQueries.GetToken(ctx, tokentable.GetTokenParams{
		UserID: validUser.ID,
		Token:  req.Token,
	})
	if err != nil {
		s.logger.WithError(err).Error("failed to get validate token")
		return fmt.Errorf("failed get validate token : %v", err)
	} else if err == sql.ErrNoRows {
		s.logger.WithError(err).Error("token is invalid")
		return fmt.Errorf("token didnt exists please resend again")
	}

	if validToken.ExpiredAt.Before(time.Now().UTC()) {
		s.logger.WithError(err).Error("token has expired")
		return fmt.Errorf("token has expired, please resend new token")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(validUser.Password), []byte(req.Password)); err == nil {
		s.logger.WithError(err).Error("failed to confirm passwword")
		return fmt.Errorf("dont use same password")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.WithError(err).Error("failed to hash password")
		return fmt.Errorf("failed to hash password : %v", err)
	}

	if err := userQueries.UpdatePassword(ctx, user.UpdatePasswordParams{
		Password: string(hash),
		ID:       validUser.ID,
	}); err != nil {
		s.logger.WithError(err).Error("failed to update password")
		return fmt.Errorf("failed to updated password : %v", err)
	}

	if err := tokenQueries.DeleteToken(ctx, tokentable.DeleteTokenParams{
		UserID:    validUser.ID,
		TokenType: "password_reset",
	}); err != nil {
		s.logger.WithError(err).Error("failed to delete token")
		return fmt.Errorf("failed to delete token : %v", err)
	}

	s.logger.Info(fmt.Sprintf("user %v success update password", validUser.ID))
	return nil
}
