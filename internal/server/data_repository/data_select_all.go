package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/model"
)

func (r *PostgresDataRepository) SelectAll(ctx context.Context, userID, dataType string) ([]model.Data, error) {
	rows, err := r.postgresPool.DB.Query(ctx,
		`
			select
			    id, owner_id, type, data, metadata, created_at, updated_at 
			from privatekeeper.data
			where owner_id = $1 and type = $2;
			`,
		userID, dataType)
	if err != nil {
		return nil, fmt.Errorf("make query: %w", err)
	}

	cards, err := pgx.CollectRows(rows, pgx.RowToStructByPos[model.Data])
	if err != nil {
		return nil, fmt.Errorf("collect row: %w", err)
	}

	return cards, nil
}
