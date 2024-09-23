package repository

import "github.com/DenisKhanov/PrivateKeeperV2/internal/server/storage/postgresql"

type PostgresDataRepository struct {
	postgresPool *postgresql.PostgresPool
}

func New(postgresPool *postgresql.PostgresPool) *PostgresDataRepository {
	return &PostgresDataRepository{postgresPool: postgresPool}
}
