package userrepository

import (
	"context"
	"database/sql"

	"github.com/ArdiSasongko/go-forum-backend/internal/sqlc/usersession"
)

type userSessionRepository struct {
	queries *usersession.Queries
}

func NewUserSessionRepository(database *sql.DB) *userSessionRepository {
	return &userSessionRepository{queries: usersession.New(database)}
}

type UserSession interface {
	InsertToken(ctx context.Context, tx *sql.Tx, model usersession.UserSession) (usersession.InsertTokenRow, error)
	GetToken(ctx context.Context, tx *sql.Tx, userid int32) (usersession.UserSession, error)
	UpdateToken(ctx context.Context, tx *sql.Tx, model usersession.UserSession) error
	DeleteToken(ctx context.Context, tx *sql.Tx, userid int32) error
}

func (r *userSessionRepository) InsertToken(ctx context.Context, tx *sql.Tx, model usersession.UserSession) (usersession.InsertTokenRow, error) {
	q := r.queries.WithTx(tx)
	return q.InsertToken(ctx, usersession.InsertTokenParams(model))
}

func (r *userSessionRepository) GetToken(ctx context.Context, tx *sql.Tx, userid int32) (usersession.UserSession, error) {
	q := r.queries.WithTx(tx)
	return q.GetToken(ctx, userid)
}

func (r *userSessionRepository) UpdateToken(ctx context.Context, tx *sql.Tx, model usersession.UserSession) error {
	q := r.queries.WithTx(tx)
	return q.UpdateToken(ctx, usersession.UpdateTokenParams{
		Token:               model.Token,
		TokenExpired:        model.TokenExpired,
		RefreshToken:        model.RefreshToken,
		RefreshTokenExpired: model.RefreshTokenExpired,
		UserID:              model.UserID,
	})
}

func (r *userSessionRepository) DeleteToken(ctx context.Context, tx *sql.Tx, userid int32) error {
	q := r.queries.WithTx(tx)
	return q.DeleteToken(ctx, userid)
}

func (r *userSessionRepository) GetTokenByToken(ctx context.Context, token string) (usersession.UserSession, error) {
	return r.queries.GetTokenByToken(ctx, token)
}
