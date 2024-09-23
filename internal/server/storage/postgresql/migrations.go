package postgresql

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/jackc/pgx/v5/stdlib"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS //

type Migrations struct {
	db *sql.DB
}

func NewMigrations(db *PostgresPool) (*Migrations, error) {
	err := goose.SetDialect("postgres")
	if err != nil {
		return nil, fmt.Errorf("goose.SetDialect: %w", err)
	}
	goose.SetBaseFS(embedMigrations)

	return &Migrations{db: stdlib.OpenDBFromPool(db.DB)}, nil
}

func (m *Migrations) Up() error {
	return goose.Up(m.db, "migrations")
}
