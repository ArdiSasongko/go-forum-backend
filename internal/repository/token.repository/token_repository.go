package tokenrepository

import (
	"context"
	"database/sql"

	tokentable "github.com/ArdiSasongko/go-forum-backend/internal/db/token"
)

type tokenRepository struct {
	queries *tokentable.Queries
}

func NewTokenRepository(database *sql.DB) *tokenRepository {
	return &tokenRepository{queries: tokentable.New(database)}
}

type TokenRepository interface {
	InsertToken(ctx context.Context, tx *sql.Tx, model tokentable.CreateTokenParams) error
	UpdateToken(ctx context.Context, tx *sql.Tx, model tokentable.UpdateTokenParams) error
	GetToken(ctx context.Context, tx *sql.Tx, userID, token int32) (tokentable.Token, error)
	DeleteToken(ctx context.Context, tx *sql.Tx, userID int32, tokenType string) error
}

func (r *tokenRepository) InsertToken(ctx context.Context, tx *sql.Tx, model tokentable.CreateTokenParams) error {
	q := r.queries.WithTx(tx)
	return q.CreateToken(ctx, tokentable.CreateTokenParams{
		UserID:    model.UserID,
		TokenType: model.TokenType,
		Token:     model.Token,
		ExpiredAt: model.ExpiredAt,
	})
}

func (r *tokenRepository) UpdateToken(ctx context.Context, tx *sql.Tx, model tokentable.UpdateTokenParams) error {
	q := r.queries.WithTx(tx)
	return q.UpdateToken(ctx, tokentable.UpdateTokenParams{
		Token:     model.Token,
		UserID:    model.UserID,
		ExpiredAt: model.ExpiredAt,
	})
}

func (r *tokenRepository) GetToken(ctx context.Context, tx *sql.Tx, userID, token int32) (tokentable.Token, error) {
	q := r.queries.WithTx(tx)
	return q.GetToken(ctx, tokentable.GetTokenParams{
		UserID: userID,
		Token:  token,
	})
}

func (r *tokenRepository) DeleteToken(ctx context.Context, tx *sql.Tx, userID int32, tokenType string) error {
	q := r.queries.WithTx(tx)
	return q.DeleteToken(ctx, tokentable.DeleteTokenParams{
		UserID:    userID,
		TokenType: tokentable.TokenType(tokenType),
	})
}
