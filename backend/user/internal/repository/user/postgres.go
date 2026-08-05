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
	userUUID, err := r.getQueries(ctx).AddUser(ctx, sqlcuser.AddUserParams{
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
	return userUUID.String(), nil
}

func (r *userRepository) GetUserByID(ctx context.Context, userID string) (*entity.User, error) {
	userUUID, err := getUUID(userID)
	if err != nil {
		return nil, fmt.Errorf("get user by id: parse user id: %w", err)
	}

	user, err := r.getQueries(ctx).GetUserByID(ctx, userUUID)

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

func (r *userRepository) GetUsersByIDs(ctx context.Context, ids []string) ([]entity.User, error) {
	uuids := make([]pgtype.UUID, len(ids))
	for i, id := range ids {
		userUUID, err := getUUID(id)
		if err != nil {
			return nil, fmt.Errorf("get users by ids: %w", err)
		}
		uuids[i] = userUUID
	}

	rawUsers, err := r.getQueries(ctx).GetUsersByIDs(ctx, uuids)
	if err != nil {
		return nil, fmt.Errorf("get users by ids: %w", err)
	}

	users := make([]entity.User, len(rawUsers))
	for i, u := range rawUsers {
		users[i] = entity.User{
			ID:       u.ID.String(),
			Username: u.Username,
		}
	}
	return users, nil
}

func (r *userRepository) UpdateCoins(
	ctx context.Context,
	userID string,
	deltaCoins int64,
) (*entity.User, error) {
	userUUID, err := getUUID(userID)
	if err != nil {
		return nil, fmt.Errorf("update coins: parse user id: %w", err)
	}

	coins, err := r.getQueries(ctx).UpdateCoins(ctx,
		sqlcuser.UpdateCoinsParams{
			UserID:     userUUID,
			DeltaCoins: deltaCoins,
		},
	)
	if err == nil {
		return &entity.User{
			ID:    userID,
			Coins: uint64(coins),
		}, nil
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("update coins: %w", err)
	}

	exists, err := r.getQueries(ctx).UserExists(ctx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("update coins: check user existence: %w", err)
	}

	if !exists {
		return nil, entity.ErrUserNotFound
	}

	return nil, entity.ErrInsufficientCoins
}

func (r *userRepository) getQueries(ctx context.Context) *sqlcuser.Queries {
	tx, err := db.ExtractTx(ctx)
	if err != nil {
		return r.queries
	}

	return sqlcuser.New(tx)
}

func getUUID(userID string) (pgtype.UUID, error) {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return pgtype.UUID{}, fmt.Errorf("get uuid: %w", err)
	}
	return pgtype.UUID{
		Bytes: parsedUserID,
		Valid: true,
	}, nil
}
