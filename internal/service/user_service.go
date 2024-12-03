package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ArdiSasongko/go-forum-backend/internal/db/user"
	"github.com/ArdiSasongko/go-forum-backend/internal/db/usersession"
	"github.com/ArdiSasongko/go-forum-backend/internal/model"
	"github.com/ArdiSasongko/go-forum-backend/internal/repository"
	"github.com/ArdiSasongko/go-forum-backend/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo        repository.UserRepository
	sessionRepo repository.UserSession
	db          *sql.DB
}

func NewUserService(repo repository.UserRepository, sessionRepo repository.UserSession, db *sql.DB) *userService {
	return &userService{
		repo:        repo,
		sessionRepo: sessionRepo,
		db:          db,
	}
}

type UserService interface {
	CreateUser(ctx context.Context, model model.UserModel) error
	LoginUser(ctx context.Context, req model.LoginRequest) (*model.ResponseLogin, error)
	RefreshToken(ctx context.Context, req model.PayloadToken, token model.RefreshToken) (string, error)
}

func (s *userService) CreateUser(ctx context.Context, req model.UserModel) error {
	tx, err := s.db.BeginTx(ctx, nil)
	defer utils.Tx(tx, err)

	// checking username in database
	_, err = s.repo.GetUser(ctx, tx, 0, req.Username, "")
	if err == nil {
		return fmt.Errorf("username already used")
	} else if err != sql.ErrNoRows {
		logrus.WithField("get username", err.Error()).Error(err.Error())
		return fmt.Errorf("failed to get user :%v", err)
	}

	// checking email in database
	_, err = s.repo.GetUser(ctx, tx, 0, "", req.Email)
	if err == nil {
		return fmt.Errorf("email already used")
	} else if err != sql.ErrNoRows {
		logrus.WithField("get email", err.Error()).Error(err.Error())
		return fmt.Errorf("failed to get user :%v", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logrus.WithField("hash password", err.Error()).Error(err.Error())
		return fmt.Errorf("failed to generated hash password :%v", err)
	}

	model := user.User{
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Password: string(hash),
		Role:     "user",
		IsValid:  sql.NullBool{Bool: false, Valid: true},
	}

	err = s.repo.CreateUser(ctx, tx, model)
	if err != nil {
		logrus.WithField("create user", err.Error()).Error(err.Error())
		return fmt.Errorf("failed to create new user :%v", err)
	}

	return nil
}

func (s *userService) LoginUser(ctx context.Context, req model.LoginRequest) (*model.ResponseLogin, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	defer utils.Tx(tx, err)

	// get user
	user, err := s.repo.GetUser(ctx, tx, 0, "", req.Email)
	if err == sql.ErrNoRows {
		logrus.WithField("get user", "user didn't exist").Error("user didn't exist")
		return nil, fmt.Errorf("invalid credentials")
	} else if err != nil {
		logrus.WithField("get user", err.Error()).Error(err.Error())
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	// check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		logrus.WithField("compare password", err.Error()).Error(err.Error())
		return nil, fmt.Errorf("invalid credentials")
	}

	// generate token and refresh token
	claims := utils.ClaimsToken{
		Username: user.Username,
		Email:    user.Email,
		Role:     string(user.Role),
		IsValid:  user.IsValid.Bool,
	}

	logrus.Info(user.IsValid.Bool, user.IsValid.Valid)
	token, err := utils.GenerateToken(ctx, claims, "token")
	if err != nil {
		logrus.WithField("generated token", err.Error()).Error(err.Error())
		return nil, fmt.Errorf("failed to generate token: %v", err)
	}

	refreshToken, err := utils.GenerateToken(ctx, claims, "refresh_token")
	if err != nil {
		logrus.WithField("generated refresh token", err.Error()).Error(err.Error())
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

	validToken, err := s.sessionRepo.GetToken(ctx, tx, user.ID)
	if err == sql.ErrNoRows {
		_, err := s.sessionRepo.InsertToken(ctx, tx, tokenModel)
		if err != nil {
			logrus.WithField("insert token", err.Error()).Error(err.Error())
			return nil, fmt.Errorf("failed to insert token: %v", err)
		}
		return &model.ResponseLogin{
			Token:        tokenModel.Token,
			RefreshToken: tokenModel.RefreshToken,
		}, nil
	} else if err != nil {
		logrus.WithField("get token", err.Error()).Error(err.Error())
		return nil, fmt.Errorf("failed to get token: %v", err)
	}

	logrus.Info(validToken.TokenExpired, currentTime)
	if validToken.RefreshTokenExpired.Before(currentTime) {
		err := s.sessionRepo.UpdateToken(ctx, tx, tokenModel)
		if err != nil {
			logrus.WithField("update token", err.Error()).Error(err.Error())
			return nil, fmt.Errorf("failed to update token: %v", err)
		}

		return &model.ResponseLogin{
			Token:        tokenModel.Token,
			RefreshToken: tokenModel.RefreshToken,
		}, nil
	} else if validToken.TokenExpired.Before(currentTime) {
		tokenModel.RefreshToken = validToken.RefreshToken
		tokenModel.RefreshTokenExpired = validToken.RefreshTokenExpired
		err := s.sessionRepo.UpdateToken(ctx, tx, tokenModel)
		if err != nil {
			logrus.WithField("update token", err.Error()).Error(err.Error())
			return nil, fmt.Errorf("failed to update token: %v", err)
		}

		return &model.ResponseLogin{
			Token:        tokenModel.Token,
			RefreshToken: validToken.RefreshToken,
		}, nil
	}

	return &model.ResponseLogin{
		Token:        validToken.Token,
		RefreshToken: validToken.RefreshToken,
	}, nil
}

func (s *userService) RefreshToken(ctx context.Context, req model.PayloadToken, token model.RefreshToken) (string, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	defer utils.Tx(tx, err)

	// get user
	user, err := s.repo.GetUser(ctx, tx, 0, req.Username, "")
	if err != nil || err == sql.ErrNoRows {
		logrus.WithField("get user", err.Error()).Error(err.Error())
		return "", fmt.Errorf("failed get user : %v", err)
	}

	// get token
	validToken, err := s.sessionRepo.GetToken(ctx, tx, user.ID)
	if err != nil || err == sql.ErrNoRows {
		logrus.WithField("get token", err.Error()).Error(err.Error())
		return "", fmt.Errorf("failed get token : %v", err)
	}

	if validToken.RefreshToken != token.Token {
		logrus.WithField("get token", "token is invalid").Error("token is invalid")
		return "", fmt.Errorf("token is invalid")
	}

	currentTime := time.Now().UTC()
	if validToken.RefreshTokenExpired.Before(currentTime) {
		logrus.WithField("get token", "token has expired").Error("token has expired")
		return "", fmt.Errorf("token has expired, please login again")
	}

	// create new token
	claims := utils.ClaimsToken{
		Username: req.Username,
		Email:    req.Email,
		Role:     req.Role,
	}

	newToken, err := utils.GenerateToken(ctx, claims, "token")
	if err != nil {
		logrus.WithField("create token", err.Error()).Error(err.Error())
		return "", fmt.Errorf("failed create token : %v", err)
	}

	model := usersession.UserSession{
		UserID:              user.ID,
		Token:               newToken,
		TokenExpired:        currentTime.Add(utils.MapToken["token"]),
		RefreshToken:        validToken.RefreshToken,
		RefreshTokenExpired: validToken.RefreshTokenExpired,
	}

	if err := s.sessionRepo.UpdateToken(ctx, tx, model); err != nil {
		logrus.WithField("update token", err.Error()).Error(err.Error())
		return "", fmt.Errorf("failed to update token : %v", err)
	}

	return newToken, nil
}
