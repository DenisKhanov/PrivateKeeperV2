package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/user/cerrors"

	"github.com/jackc/pgx/v5"

	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/model"
)

func (r *PostgresUserRepository) SelectByLogin(ctx context.Context, login string) (model.User, error) {
	rows, err := r.postgresPool.DB.Query(ctx,
		`
			select 
				id, login, password, crypt_key, created_at, updated_at
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
