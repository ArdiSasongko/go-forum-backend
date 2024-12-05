package userhandler

import (
	"database/sql"

	userservice "github.com/ArdiSasongko/go-forum-backend/internal/service/user.service"
	imageuser "github.com/ArdiSasongko/go-forum-backend/internal/sqlc/image_user"
	tokentable "github.com/ArdiSasongko/go-forum-backend/internal/sqlc/token"
	"github.com/ArdiSasongko/go-forum-backend/internal/sqlc/user"
	"github.com/ArdiSasongko/go-forum-backend/internal/sqlc/usersession"
)

type userHandler struct {
	service userservice.UserService
}

func NewUserHandler(service userservice.UserService) *userHandler {
	return &userHandler{service: service}
}

var db *sql.DB
var queries = userservice.Queries{
	UserQueries:        user.New(db),
	TokenQueries:       tokentable.New(db),
	ImageUserQueries:   imageuser.New(db),
	UserSessionQueries: usersession.New(db),
}
