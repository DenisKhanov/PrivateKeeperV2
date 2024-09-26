package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/model"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/storage/postgresql"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/user/cerrors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type PostgresUserRepository struct {
	postgresPool *postgresql.PostgresPool
}

func New(postgresPool *postgresql.PostgresPool) *PostgresUserRepository {
	return &PostgresUserRepository{postgresPool: postgresPool}
}

func (r *PostgresUserRepository) SelectByLogin(ctx context.Context, login string) (model.User, error) {
	rows, err := r.postgresPool.DB.Query(ctx,
		`
			select 
				id, login, password, crypt_key, created_at
			from privatekeeper.user
			where login = $1;
			`,
		login)
	if err != nil {
		return model.User{}, fmt.Errorf("insert user: %w", err)
	}

	savedUser, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[model.User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, cerrors.ErrUserNotFound
		}

		return model.User{}, fmt.Errorf("insert user: %w", err)
	}

	return savedUser, nil
}

func (r *PostgresUserRepository) SelectKeyByID(ctx context.Context, userID string) ([]byte, error) {
	var userKey []byte
	err := r.postgresPool.DB.QueryRow(ctx,
		`
			select 
				crypt_key
			from privatekeeper.user
			where id = $1;
			`,
		userID).Scan(&userKey)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	return userKey, nil
}

func (r *PostgresUserRepository) Insert(ctx context.Context, user model.User) (model.User, error) {
	rows, err := r.postgresPool.DB.Query(ctx,
		`
			insert into privatekeeper.user
				(id, login, password, crypt_key, created_at)
			values
				($1, $2, $3, $4, NOW(), NOW())
			returning id, login, password, crypt_key, created_at;
			`,
		user.ID,
		user.Login,
		user.Password,
		user.CryptKey)
	if err != nil {
		return model.User{}, fmt.Errorf("make query: %w", err)
	}

	savedUser, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[model.User])
	var e *pgconn.PgError
	if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation {
		return model.User{}, fmt.Errorf("collect row: %w", cerrors.ErrUserAlreadyExists)
	}

	if err != nil {
		return model.User{}, fmt.Errorf("collect row: %w", err)
	}

	return savedUser, nil
}
