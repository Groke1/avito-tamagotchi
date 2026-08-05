package token

import (
	"context"
	"errors"
	"fmt"

	"github.com/cayman444/avito-gamification-hackathon.pkg/db"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
	sqlctoken "github.com/cayman444/avito-gamification-hackathon.user/internal/repository/token/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type tokenRepository struct {
	db      sqlctoken.DBTX
	queries *sqlctoken.Queries
}

func NewTokenRepository(qdb sqlctoken.DBTX) *tokenRepository {
	return &tokenRepository{
		db:      qdb,
		queries: sqlctoken.New(qdb),
	}
}

func (r *tokenRepository) AddToken(ctx context.Context, userID string, token entity.RefreshToken) error {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("add refresh token: parse user id: %w", err)
	}

	err = r.getQueries(ctx).AddToken(ctx, sqlctoken.AddTokenParams{
		UserID: pgtype.UUID{
			Bytes: parsedUserID,
			Valid: true,
		},
		TokenHash: token.TokenHash,
		ExpiresAt: pgtype.Timestamptz{
			Time:  token.ExpiresAt,
			Valid: true,
		},
	})
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
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

func (r *tokenRepository) getQueries(ctx context.Context) *sqlctoken.Queries {
	tx, err := db.ExtractTx(ctx)
	if err != nil {
		return r.queries
	}

	return sqlctoken.New(tx)
}
