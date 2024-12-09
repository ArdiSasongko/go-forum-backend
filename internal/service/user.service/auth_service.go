package userservice

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ArdiSasongko/go-forum-backend/env"
	"github.com/ArdiSasongko/go-forum-backend/internal/model"
	imageuser "github.com/ArdiSasongko/go-forum-backend/internal/sqlc/image_user"
	tokentable "github.com/ArdiSasongko/go-forum-backend/internal/sqlc/token"
	"github.com/ArdiSasongko/go-forum-backend/internal/sqlc/user"
	"github.com/ArdiSasongko/go-forum-backend/internal/sqlc/usersession"
	cld "github.com/ArdiSasongko/go-forum-backend/pkg/cloudinary"
	"github.com/ArdiSasongko/go-forum-backend/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func (s *userService) CreateUser(ctx context.Context, queries Queries, req model.UserModel) error {
	tx, err := s.db.BeginTx(ctx, nil)
	defer utils.Tx(tx, err)
	userQueries := queries.UserQueries.WithTx(tx)
	tokenQueries := queries.TokenQueries.WithTx(tx)
	imageQueries := queries.ImageUserQueries.WithTx(tx)
	// checking username in database
	_, err = userQueries.GetUser(ctx, user.GetUserParams{
		ID:       0,
		Username: req.Username,
		Email:    "",
	})
	if err == nil {
		return fmt.Errorf("username already used")
	} else if err != sql.ErrNoRows {
		s.logger.WithError(err).Error("failed to get username")
		return fmt.Errorf("failed to get user :%v", err)
	}

	// checking email in database
	_, err = userQueries.GetUser(ctx, user.GetUserParams{
		ID:       0,
		Username: "",
		Email:    req.Email,
	})
	if err == nil {
		return fmt.Errorf("email already used")
	} else if err != sql.ErrNoRows {
		s.logger.WithError(err).Error("failed to get user")
		return fmt.Errorf("failed to get user :%v", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.WithError(err).Error("failed to hash password")
		return fmt.Errorf("failed to generated hash password :%v", err)
	}

	id, err := userQueries.CreateUser(ctx, user.CreateUserParams{
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Password: string(hash),
		Role:     "user",
		IsValid:  sql.NullBool{Bool: false, Valid: true},
	})
	if err != nil {
		s.logger.WithError(err).Error("failed to create new user")
		return fmt.Errorf("failed to create new user :%v", err)
	}

	// upload image
	url := env.GetEnv("CLOUDINARY_URL", "")
	profile, _, err := utils.GetProfile()
	if err != nil {
		s.logger.WithError(err).Error("failed to get default profile image")
		return fmt.Errorf("failed to get profile :%v", err)
	}

	imageUrl, publicID, err := cld.UploadImageByte(ctx, profile, url, "forum-profile")
	if err != nil {
		s.logger.WithError(err).Error("failed to upload image profile")
		return fmt.Errorf("failed to upload image profile :%v", err)
	}

	// insert image
	if err := imageQueries.CreateImage(ctx, imageuser.CreateImageParams{
		UserID:   id,
		ImageUrl: imageUrl,
	}); err != nil {
		if err := cld.DestroyImage(ctx, url, publicID); err != nil {
			s.logger.WithError(err).Error("failed to delete image profile")
			return fmt.Errorf("failed to delete image profile :%v", err)
		}
		s.logger.WithError(err).Error("failed to insert image")
		return fmt.Errorf("failed to insert image profile :%v", err)
	}

	// create token
	validationToken := utils.GenToken()
	err = tokenQueries.CreateToken(ctx, tokentable.CreateTokenParams{
		UserID:    id,
		TokenType: "email",
		Token:     int32(validationToken),
		ExpiredAt: time.Now().UTC().Add(5 * time.Minute),
	})
	if err != nil {
		s.logger.WithError(err).Error("failed to create new token")
		return fmt.Errorf("failed to create new token :%v", err)
	}

	logrus.Info(validationToken)
	err = utils.SendToken(req.Email, "email", int32(validationToken))
	if err != nil {
		s.logger.WithError(err).Error("failed to send email")
		return fmt.Errorf("failed to send email :%v", err)
	}

	s.logger.Info("success register")
	return nil
}

func (s *userService) LoginUser(ctx context.Context, queries Queries, req model.LoginRequest) (*model.ResponseLogin, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	defer utils.Tx(tx, err)
	userQueries := queries.UserQueries.WithTx(tx)
	userSessionQueries := queries.UserSessionQueries.WithTx(tx)
	// get user
	user, err := userQueries.GetUser(ctx, user.GetUserParams{
		ID:       0,
		Username: "",
		Email:    req.Email,
	})
	if err == sql.ErrNoRows {
		s.logger.WithError(err).Error("user didnt exists")
		return nil, fmt.Errorf("invalid credentials")
	} else if err != nil {
		s.logger.WithError(err).Error("failed to get user")
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	// check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		s.logger.WithError(err).Error("failed to confirm password")
		return nil, fmt.Errorf("invalid credentials")
	}

	// generate token and refresh token
	claims := utils.ClaimsToken{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     string(user.Role),
		IsValid:  user.IsValid.Bool,
	}

	logrus.Info(user.IsValid.Bool, user.IsValid.Valid)
	token, err := utils.GenerateToken(ctx, claims, "token")
	if err != nil {
		s.logger.WithError(err).Error("failed to generate token")
		return nil, fmt.Errorf("failed to generate token: %v", err)
	}

	refreshToken, err := utils.GenerateToken(ctx, claims, "refresh_token")
	if err != nil {
		s.logger.WithError(err).Error("failed to generate refresh token")
		return nil, fmt.Errorf("failed to generate refresh token: %v", err)
	}

	currentTime := time.Now().UTC()
	tokenModel := usersession.UserSession{
		UserID:              user.ID,
		Token:               token,
		TokenExpired:        currentTime.Add(utils.MapToken["token"]),
		RefreshToken:        refreshToken,
		RefreshTokenExpired: currentTime.Add(utils.MapToken["refresh_token"]),
	}

	validToken, err := userSessionQueries.GetToken(ctx, user.ID)
	if err == sql.ErrNoRows {
		_, err := userSessionQueries.InsertToken(ctx, usersession.InsertTokenParams(tokenModel))
		if err != nil {
			s.logger.WithError(err).Error("failed to insert token")
			return nil, fmt.Errorf("failed to insert token: %v", err)
		}
		return &model.ResponseLogin{
			Token:        tokenModel.Token,
			RefreshToken: tokenModel.RefreshToken,
		}, nil
	} else if err != nil {
		s.logger.WithError(err).Error("failed to get token")
		return nil, fmt.Errorf("failed to get token: %v", err)
	}

	logrus.Info(validToken.TokenExpired, currentTime)
	if validToken.RefreshTokenExpired.Before(currentTime) {
		err := userSessionQueries.UpdateToken(ctx, usersession.UpdateTokenParams{
			UserID:              user.ID,
			Token:               token,
			TokenExpired:        currentTime.Add(utils.MapToken["token"]),
			RefreshToken:        refreshToken,
			RefreshTokenExpired: currentTime.Add(utils.MapToken["refresh_token"]),
		})
		if err != nil {
			s.logger.WithError(err).Error("failed to update token")
			return nil, fmt.Errorf("failed to update token: %v", err)
		}

		return &model.ResponseLogin{
			Token:        tokenModel.Token,
			RefreshToken: tokenModel.RefreshToken,
		}, nil
	} else if validToken.TokenExpired.Before(currentTime) {
		err := userSessionQueries.UpdateToken(ctx, usersession.UpdateTokenParams{
			UserID:              user.ID,
			Token:               token,
			TokenExpired:        currentTime.Add(utils.MapToken["token"]),
			RefreshToken:        validToken.RefreshToken,
			RefreshTokenExpired: validToken.RefreshTokenExpired,
		})
		if err != nil {
			s.logger.WithError(err).Error("failed to update token")
			return nil, fmt.Errorf("failed to update token: %v", err)
		}

		return &model.ResponseLogin{
			Token:        tokenModel.Token,
			RefreshToken: validToken.RefreshToken,
		}, nil
	}

	s.logger.Info(fmt.Sprintf("user %v success login", user.ID))
	return &model.ResponseLogin{
		Token:        validToken.Token,
		RefreshToken: validToken.RefreshToken,
	}, nil
}

func (s *userService) RefreshToken(ctx context.Context, queries Queries, req model.PayloadToken, token model.RefreshToken) (string, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	defer utils.Tx(tx, err)
	userQueris := queries.UserQueries.WithTx(tx)
	userSessionQueries := queries.UserSessionQueries.WithTx(tx)
	// get user
	user, err := userQueris.GetUser(ctx, user.GetUserParams{
		ID:       0,
		Username: req.Username,
		Email:    "",
	})
	if err != nil || err == sql.ErrNoRows {
		s.logger.WithError(err).Error("failed to get user")
		return "", fmt.Errorf("failed get user : %v", err)
	}

	// get token
	validToken, err := userSessionQueries.GetToken(ctx, user.ID)
	if err != nil || err == sql.ErrNoRows {
		s.logger.WithError(err).Error("failed to get token")
		return "", fmt.Errorf("failed get token : %v", err)
	}

	if validToken.RefreshToken != token.Token {
		s.logger.WithError(err).Error("token is invalid")
		return "", fmt.Errorf("token is invalid")
	}

	currentTime := time.Now().UTC()
	if validToken.RefreshTokenExpired.Before(currentTime) {
		s.logger.WithError(err).Error("token has expired")
		return "", fmt.Errorf("token has expired, please login again")
	}

	// create new token
	claims := utils.ClaimsToken{
		UserID:   req.UserID,
		Username: req.Username,
		Email:    req.Email,
		Role:     req.Role,
	}

	newToken, err := utils.GenerateToken(ctx, claims, "token")
	if err != nil {
		s.logger.WithError(err).Error("failed to create token")
		return "", fmt.Errorf("failed create token : %v", err)
	}

	if err := userSessionQueries.UpdateToken(ctx, usersession.UpdateTokenParams{
		UserID:              user.ID,
		Token:               newToken,
		TokenExpired:        currentTime.Add(utils.MapToken["token"]),
		RefreshToken:        validToken.RefreshToken,
		RefreshTokenExpired: validToken.RefreshTokenExpired,
	}); err != nil {
		s.logger.WithError(err).Error("failed to update token")
		return "", fmt.Errorf("failed to update token : %v", err)
	}

	s.logger.Info(fmt.Sprintf("user %v success get refresh token", user.ID))
	return newToken, nil
}

func (s *userService) Logout(ctx context.Context, queries Queries, id int32) error {
	tx, err := s.db.BeginTx(ctx, nil)
	defer utils.Tx(tx, err)
	userQueries := queries.UserQueries.WithTx(tx)
	userSessionQueries := queries.UserSessionQueries.WithTx(tx)

	_, err = userQueries.GetUser(ctx, user.GetUserParams{
		ID:       id,
		Username: "",
		Email:    "",
	})

	if err == sql.ErrNoRows {
		s.logger.WithError(err).Error("failed to get user")
		return fmt.Errorf("failed to get user : %v", err.Error())
	} else if err != nil {
		s.logger.WithError(err).Error("failed to get user")
		return fmt.Errorf("failed to get user : %v", err.Error())
	}

	if err := userSessionQueries.DeleteToken(ctx, id); err != nil {
		s.logger.WithError(err).Error("failed to delete user session")
		return fmt.Errorf("failed to delete user session : %v", err.Error())
	}

	s.logger.Info(fmt.Sprintf("user %v success logout", id))
	return nil
}
