package database

import (
	"context"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBStorage struct {
	DB *pgxpool.Pool
}

func New(ctx context.Context, addr string) (*DBStorage, error) {
	pool, err := pgxpool.New(ctx, addr)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}
	return &DBStorage{DB: pool}, nil
}

func (db *DBStorage) Close() {
	db.DB.Close()
}

func ApplyMigrations(addr string) error {
	migrationPath := "file://migrations"

	m, err := migrate.New(migrationPath, addr)
	if err != nil {
		return err
	}

	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
