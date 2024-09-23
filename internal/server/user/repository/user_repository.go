package repository

import "github.com/DenisKhanov/PrivateKeeperV2/internal/server/storage/postgresql"

type PostgresUserRepository struct {
	postgresPool *postgresql.PostgresPool
}

func New(postgresPool *postgresql.PostgresPool) *PostgresUserRepository {
	return &PostgresUserRepository{postgresPool: postgresPool}
}
