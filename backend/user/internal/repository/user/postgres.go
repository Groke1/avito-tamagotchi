package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/cayman444/avito-gamification-hackathon.pkg/db"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
	sqlcuser "github.com/cayman444/avito-gamification-hackathon.user/internal/repository/user/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type userRepository struct {
	db      sqlcuser.DBTX
	queries *sqlcuser.Queries
}

func NewUserRepository(qdb sqlcuser.DBTX) *userRepository {
	return &userRepository{
		db:      qdb,
		queries: sqlcuser.New(qdb),
	}
}

func (r *userRepository) AddUser(ctx context.Context, user entity.User) (userID string, err error) {
	uuid, err := r.getQueries(ctx).AddUser(ctx, sqlcuser.AddUserParams{
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Coins:        int64(user.Coins),
	})

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "users_username_unique", "users_username_key":
				return "", entity.ErrUsernameAlreadyExists
			case "users_email_unique", "users_email_key":
				return "", entity.ErrEmailAlreadyExists
			default:
				return "", fmt.Errorf(
					"add user: unique constraint %q violated: %w",
					pgErr.ConstraintName,
					err,
				)
			}
		}
		return "", fmt.Errorf("add user: %w", err)
	}
	return uuid.String(), nil
}

func (r *userRepository) GetUserByID(ctx context.Context, userID string) (*entity.User, error) {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("get user by id: parse user id: %w", err)
	}

	user, err := r.getQueries(ctx).GetUserByID(
		ctx,
		pgtype.UUID{
			Bytes: parsedUserID,
			Valid: true,
		})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.ErrUserNotFound
		}

		return nil, fmt.Errorf("get user by id: %w", err)
	}

	if user.Coins < 0 {
		return nil, fmt.Errorf("get user by id: invalid negative coins value: %d", user.Coins)
	}

	return &entity.User{
		ID:       user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
		Coins:    uint64(user.Coins),
	}, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := r.getQueries(ctx).GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.ErrUserNotFound
		}

		return nil, fmt.Errorf("get user by email: %w", err)
	}

	if user.Coins < 0 {
		return nil, fmt.Errorf("get user by email: invalid negative coins value: %d", user.Coins)
	}

	return &entity.User{
		ID:           user.ID.String(),
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Coins:        uint64(user.Coins),
	}, nil
}

func (r *userRepository) getQueries(ctx context.Context) *sqlcuser.Queries {
	tx, err := db.ExtractTx(ctx)
	if err != nil {
		return r.queries
	}

	return sqlcuser.New(tx)
}
