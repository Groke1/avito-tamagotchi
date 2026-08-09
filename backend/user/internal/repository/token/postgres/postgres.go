package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/cayman444/avito-gamification-hackathon.pkg/db"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/repository/converter"
	sqlc "github.com/cayman444/avito-gamification-hackathon.user/internal/repository/token/postgres/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type tokenRepository struct {
	db      sqlc.DBTX
	queries *sqlc.Queries
}

func NewTokenRepository(qdb sqlc.DBTX) *tokenRepository {
	return &tokenRepository{
		db:      qdb,
		queries: sqlc.New(qdb),
	}
}

func (r *tokenRepository) AddToken(ctx context.Context, userID string, token entity.RefreshToken) error {
	userUUID, err := converter.StringToUUID(userID)
	if err != nil {
		return fmt.Errorf("add refresh token: parse user id: %w", err)
	}

	err = r.getQueries(ctx).AddToken(ctx, sqlc.AddTokenParams{
		UserID:    userUUID,
		TokenHash: token.TokenHash,
		ExpiresAt: pgtype.Timestamptz{
			Time:  token.ExpiresAt,
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
	token, err := r.getQueries(ctx).GetRefreshTokenByHashForUpdate(ctx, hash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.ErrRefreshTokenNotFound
		}

		return nil, fmt.Errorf("get refresh token by hash: %w", err)
	}

	if !token.ExpiresAt.Valid {
		return nil, errors.New("get refresh token by hash: expires_at is null")
	}

	return &entity.RefreshToken{
		UserID:    token.UserID.String(),
		TokenHash: token.TokenHash,
		ExpiresAt: token.ExpiresAt.Time,
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

	if err := r.getQueries(ctx).DeleteSession(ctx, sqlc.DeleteSessionParams{
		UserID:    userUUID,
		TokenHash: tokenHash,
	}); err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	return nil
}

func (r *tokenRepository) getQueries(ctx context.Context) *sqlc.Queries {
	tx, err := db.ExtractTx(ctx)
	if err != nil {
		return r.queries
	}

	return sqlc.New(tx)
}
