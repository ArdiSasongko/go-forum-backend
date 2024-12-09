package userservice

import (
	"context"
	"database/sql"

	"github.com/ArdiSasongko/go-forum-backend/internal/model"
	imageuser "github.com/ArdiSasongko/go-forum-backend/internal/sqlc/image_user"
	tokentable "github.com/ArdiSasongko/go-forum-backend/internal/sqlc/token"
	"github.com/ArdiSasongko/go-forum-backend/internal/sqlc/user"
	"github.com/ArdiSasongko/go-forum-backend/internal/sqlc/usersession"
	"github.com/sirupsen/logrus"
)

type userService struct {
	db     *sql.DB
	logger *logrus.Logger
}

func NewUserService(db *sql.DB, logger *logrus.Logger) *userService {
	return &userService{
		db:     db,
		logger: logger,
	}
}

type Queries struct {
	UserQueries        *user.Queries
	TokenQueries       *tokentable.Queries
	ImageUserQueries   *imageuser.Queries
	UserSessionQueries *usersession.Queries
}

type UserService interface {
	CreateUser(ctx context.Context, queries Queries, req model.UserModel) error
	LoginUser(ctx context.Context, queries Queries, req model.LoginRequest) (*model.ResponseLogin, error)
	Logout(ctx context.Context, queries Queries, id int32) error
	RefreshToken(ctx context.Context, queries Queries, req model.PayloadToken, token model.RefreshToken) (string, error)
	ValidateEmail(ctx context.Context, queries Queries, payload model.ValidatePayload) error
	ResendEmail(ctx context.Context, queries Queries, payload model.ValidatePayload) error
	ResetPassword(ctx context.Context, queries Queries, req model.SendEmail) error
	ConfirmPassword(ctx context.Context, queries Queries, req model.ResetPassword) error
	GetProfile(ctx context.Context, queris Queries, email string) (*model.ProfileModel, error)
	UpdateProfile(ctx context.Context, queries Queries, req model.UpdateProfile) error
	UpdateUser(ctx context.Context, queries Queries, req model.UpdateUser, email string) error
}
