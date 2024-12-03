package userrepository

import (
	"context"
	"database/sql"

	"github.com/ArdiSasongko/go-forum-backend/internal/sqlc/user"
)

type userRepository struct {
	queries *user.Queries
}

func NewuserRepository(database *sql.DB) *userRepository {
	return &userRepository{queries: user.New(database)}
}

type UserRepository interface {
	CreateUser(ctx context.Context, tx *sql.Tx, model user.User) (int32, error)
	GetUser(ctx context.Context, tx *sql.Tx, id int32, username, email string) (user.User, error)
	ValidateUser(ctx context.Context, tx *sql.Tx, userId int32) error
	UpdatePassword(ctx context.Context, tx *sql.Tx, password string, userID int32) error
}

func (r *userRepository) CreateUser(ctx context.Context, tx *sql.Tx, model user.User) (int32, error) {
	q := r.queries.WithTx(tx)
	return q.CreateUser(ctx, user.CreateUserParams{
		Name:     model.Name,
		Username: model.Username,
		Email:    model.Email,
		Password: model.Password,
		Role:     model.Role,
		IsValid:  model.IsValid,
	})
}

func (r *userRepository) GetUser(ctx context.Context, tx *sql.Tx, id int32, username, email string) (user.User, error) {
	q := r.queries.WithTx(tx)
	return q.GetUser(ctx, user.GetUserParams{
		ID:       id,
		Username: username,
		Email:    email,
	})
}

func (r *userRepository) ValidateUser(ctx context.Context, tx *sql.Tx, userId int32) error {
	q := r.queries.WithTx(tx)
	return q.ValidateUser(ctx, userId)
}

func (r *userRepository) UpdatePassword(ctx context.Context, tx *sql.Tx, password string, userID int32) error {
	q := r.queries.WithTx(tx)
	return q.UpdatePassword(ctx, user.UpdatePasswordParams{
		Password: password,
		ID:       userID,
	})
}
