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

// PostgresDataRepository defines a repository that interacts with PostgreSQL to manage data.
type PostgresDataRepository struct {
	postgresPool *postgresql.PostgresPool // Connection pool to PostgreSQL database
}

// New creates a new PostgresDataRepository instance with the provided PostgreSQL connection pool.
func New(postgresPool *postgresql.PostgresPool) *PostgresDataRepository {
	return &PostgresDataRepository{postgresPool: postgresPool}
}

// SelectAll retrieves all data entries of a specific type for a user from the database.
func (r *PostgresDataRepository) SelectAll(ctx context.Context, userID, dataType string) ([]model.Data, error) {
	rows, err := r.postgresPool.DB.Query(ctx,
		`
			select
			    id, owner_id, type, data, metadata, created_at
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

// Insert saves a new data entry into the database and returns the saved entry.
func (r *PostgresDataRepository) Insert(ctx context.Context, data model.Data) (model.Data, error) {
	rows, err := r.postgresPool.DB.Query(ctx,
		`
			insert into privatekeeper.data
			    (id, owner_id, type, data, metadata, created_at) 
			values
				($1, $2, $3, $4, $5, now())
			returning id, owner_id, type, data, metadata, created_at;
			`,
		data.ID,
		data.OwnerID,
		data.Type,
		data.Data,
		data.MetaData)
	if err != nil {
		return model.Data{}, fmt.Errorf("make query: %w", err)
	}

	savedData, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[model.Data])
	var e *pgconn.PgError
	if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation {
		return model.Data{}, fmt.Errorf("collect row: %w", cerrors.ErrUserAlreadyExists)
	}

	if err != nil {
		return model.Data{}, fmt.Errorf("collect row: %w", err)
	}

	return savedData, nil
}

// SelectByID retrieves a specific data entry by its ID for a user from the database.
func (r *PostgresDataRepository) SelectByID(ctx context.Context, userID, dataType, dataID string) (model.Data, error) {
	row, err := r.postgresPool.DB.Query(ctx,
		`
			select
			    id, owner_id, type, data, metadata, created_at
			from privatekeeper.data
			where owner_id = $1 and type = $2 and id = $3;
			`,
		userID, dataType, dataID)
	if err != nil {
		return model.Data{}, fmt.Errorf("make query: %w", err)
	}

	data, err := pgx.CollectOneRow(row, pgx.RowToStructByPos[model.Data])
	if err != nil {
		return model.Data{}, fmt.Errorf("collect row: %w", err)
	}

	return data, nil
}
