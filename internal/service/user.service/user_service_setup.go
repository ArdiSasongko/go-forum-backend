package userservice

import (
	"context"
	"database/sql"

	"github.com/ArdiSasongko/go-forum-backend/internal/model"
	imagesuserrepository "github.com/ArdiSasongko/go-forum-backend/internal/repository/images.user.repository"
	tokenrepository "github.com/ArdiSasongko/go-forum-backend/internal/repository/token.repository"
	userrepository "github.com/ArdiSasongko/go-forum-backend/internal/repository/user.repository"
)

type userService struct {
	repo        userrepository.UserRepository
	sessionRepo userrepository.UserSession
	tokenRepo   tokenrepository.TokenRepository
	imageRepo   imagesuserrepository.ImageUserRepository
	db          *sql.DB
}

func NewUserService(
	repo userrepository.UserRepository,
	sessionRepo userrepository.UserSession,
	tokenRepo tokenrepository.TokenRepository,
	imageRepo imagesuserrepository.ImageUserRepository,
	db *sql.DB) *userService {
	return &userService{
		repo:        repo,
		sessionRepo: sessionRepo,
		tokenRepo:   tokenRepo,
		imageRepo:   imageRepo,
		db:          db,
	}
}

type UserService interface {
	CreateUser(ctx context.Context, model model.UserModel) error
	LoginUser(ctx context.Context, req model.LoginRequest) (*model.ResponseLogin, error)
	RefreshToken(ctx context.Context, req model.PayloadToken, token model.RefreshToken) (string, error)
	ValidateEmail(ctx context.Context, payload model.ValidatePayload) error
	ResendEmail(ctx context.Context, payload model.ValidatePayload) error
	ResetPassword(ctx context.Context, req model.SendEmail) error
	ConfirmPassword(ctx context.Context, req model.ResetPassword) error
	GetProfile(ctx context.Context, username string) (*model.ProfileModel, error)
	UpdateProfile(ctx context.Context, req model.UpdateProfile) error
	UpdateUser(ctx context.Context, req model.UpdateUser, username string) error
}
