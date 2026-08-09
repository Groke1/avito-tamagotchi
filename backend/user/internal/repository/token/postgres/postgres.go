package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/cayman444/avito-gamification-hackathon.pkg/db"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/repository/converter"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/repository/token/postgres/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type tokenRepository struct {
	db      token.DBTX
	queries *token.Queries
}

func NewTokenRepository(qdb token.DBTX) *tokenRepository {
	return &tokenRepository{
		db:      qdb,
		queries: token.New(qdb),
	}
}

func (r *tokenRepository) AddToken(ctx context.Context, userID string, rToken entity.RefreshToken) error {
	userUUID, err := converter.StringToUUID(userID)
	if err != nil {
		return fmt.Errorf("add refresh token: parse user id: %w", err)
	}

	err = r.getQueries(ctx).AddToken(ctx, token.AddTokenParams{
		UserID:    userUUID,
		TokenHash: rToken.TokenHash,
		ExpiresAt: pgtype.Timestamptz{
			Time:  rToken.ExpiresAt,
			Valid: true,
		},
	})
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			switch pgErr.Code {
			case "23505":
				return fmt.Errorf("add refresh token: token hash already exists: %w", err)

			case "23503":
				return fmt.Errorf("add refresh token: user does not exist: %w", err)
			}
		}

		return fmt.Errorf("add refresh token: %w", err)
	}

	return nil
}

func (r *tokenRepository) GetRefreshTokenByHashForUpdate(ctx context.Context, hash string) (*entity.RefreshToken, error) {
	rToken, err := r.getQueries(ctx).GetRefreshTokenByHashForUpdate(ctx, hash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.ErrRefreshTokenNotFound
		}

		return nil, fmt.Errorf("get refresh token by hash: %w", err)
	}

	if !rToken.ExpiresAt.Valid {
		return nil, errors.New("get refresh token by hash: expires_at is null")
	}

	return &entity.RefreshToken{
		UserID:    rToken.UserID.String(),
		TokenHash: rToken.TokenHash,
		ExpiresAt: rToken.ExpiresAt.Time,
	}, nil
}

func (r *tokenRepository) ConsumeRefreshToken(ctx context.Context, hash string) (*entity.RefreshToken, error) {
	token, err := r.GetRefreshTokenByHashForUpdate(ctx, hash)
	if err != nil {
		return nil, err
	}

	if err := r.DeleteRefreshTokenByHash(ctx, hash); err != nil {
		return nil, err
	}

	return token, nil
}

func (r *tokenRepository) DeleteRefreshTokenByHash(ctx context.Context, hash string) error {
	if err := r.getQueries(ctx).DeleteRefreshTokenByHash(ctx, hash); err != nil {
		return fmt.Errorf("delete refresh token by hash: %w", err)
	}

	return nil
}

func (r *tokenRepository) DeleteExpiredTokens(ctx context.Context) error {
	if err := r.getQueries(ctx).DeleteExpiredTokens(ctx); err != nil {
		return fmt.Errorf("delete expired tokens: %w", err)
	}
	return nil
}

func (r *tokenRepository) DeleteSession(ctx context.Context, userID, tokenHash string) error {
	userUUID, err := converter.StringToUUID(userID)
	if err != nil {
		return fmt.Errorf("delete session: %w", err)
	}

	if err := r.getQueries(ctx).DeleteSession(ctx, token.DeleteSessionParams{
		UserID:    userUUID,
		TokenHash: tokenHash,
	}); err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	return nil
}

func (r *tokenRepository) getQueries(ctx context.Context) *token.Queries {
	tx, err := db.ExtractTx(ctx)
	if err != nil {
		return r.queries
	}

	return token.New(tx)
}
