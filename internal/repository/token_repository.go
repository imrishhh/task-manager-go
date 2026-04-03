package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	apperr "github.com/nullrish/task-manager-go/internal/errors"
	"github.com/nullrish/task-manager-go/internal/model"
)

type TokenRepository interface {
	Store(ctx context.Context, userID uuid.UUID, token string, tokenType string, expiresAt time.Time) (*model.UserToken, error)
	FindByToken(ctx context.Context, token string) (*model.UserToken, error)
	Revoke(ctx context.Context, token string) error
	RevokeByUserID(ctx context.Context, userID uuid.UUID) error
}

type tokenRepo struct {
	db *sql.DB
}

func NewTokenRepository(db *sql.DB) TokenRepository {
	return &tokenRepo{db: db}
}

func (r *tokenRepo) Store(ctx context.Context, userID uuid.UUID, token string, tokenType string, expiresAt time.Time) (*model.UserToken, error) {
	query := `
		INSERT INTO user_tokens (user_id, token, token_type, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, token, token_type, expires_at, created_at, revoked;
	`
	var resToken model.UserToken
	err := r.db.QueryRowContext(ctx, query, userID, token, tokenType, expiresAt).Scan(
		&resToken.ID,
		&resToken.UserID,
		&resToken.TokenType,
		&resToken.ExpiresAt,
		&resToken.CreatedAt,
		&resToken.Revoked,
	)
	if err != nil {
		return nil, &apperr.DatabaseError{Message: err.Error()}
	}
	return &resToken, nil
}

func (r *tokenRepo) FindByToken(ctx context.Context, token string) (*model.UserToken, error) {
	query := `
		SELECT id, user_id, token, token_type, expires_at, created_at, revoked
		FROM user_tokens WHERE token = $1
	`
	var resToken model.UserToken
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&resToken.ID,
		&resToken.UserID,
		&resToken.Token,
		&resToken.TokenType,
		&resToken.ExpiresAt,
		&resToken.CreatedAt,
		&resToken.Revoked,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &apperr.NotFoundError{Resource: "token", ID: token}
		}
		return nil, &apperr.DatabaseError{Message: err.Error()}
	}
	return &resToken, nil
}

func (r *tokenRepo) Revoke(ctx context.Context, token string) error {
	query := `
		UPDATE user_tokens SET revoked = TRUE
		WHERE token = $1
	`
	res, err := r.db.ExecContext(ctx, query, token)
	if rows, _ := res.RowsAffected(); rows == 0 {
		return &apperr.NotFoundError{Resource: "token", ID: token}
	}
	if err != nil {
		return &apperr.DatabaseError{Message: err.Error()}
	}
	return nil
}

func (r *tokenRepo) RevokeByUserID(ctx context.Context, userID uuid.UUID) error {
	query := `
		UPDATE user_tokens SET revoked = TRUE
		WHERE user_id = $1
	`
	res, err := r.db.ExecContext(ctx, query, userID)
	if rows, _ := res.RowsAffected(); rows == 0 {
		return &apperr.NotFoundError{Resource: "user_id", ID: userID.String()}
	}
	if err != nil {
		return &apperr.DatabaseError{Message: err.Error()}
	}
	return nil
}
