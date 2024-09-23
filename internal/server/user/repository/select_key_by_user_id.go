package repository

import (
	"context"
	"fmt"
)

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
