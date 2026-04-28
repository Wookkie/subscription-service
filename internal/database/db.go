package database

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBStorage struct {
	DB *pgxpool.Pool
}

func New(ctx context.Context, addr string) (*DBStorage, error) {
	log.Info().Msg("connecting to database")
	pool, err := pgxpool.New(ctx, addr)
	if err != nil {
		log.Error().Err(err).Msg("database connection failed")
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		log.Error().Err(err).Msg("database ping failed")
		return nil, err
	}

	log.Info().Msg("database connected successfully")
	return &DBStorage{DB: pool}, nil
}

func (db *DBStorage) Close() {
	log.Info().Msg("closing database connection")
	db.DB.Close()
}

func ApplyMigrations(addr string) error {
	log.Info().Msg("applying migrations")
	migrationPath := "file://migrations"

	m, err := migrate.New(migrationPath, addr)
	if err != nil {
		log.Error().Err(err).Msg("migration init failed")
		return err
	}

	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Error().Err(err).Msg("migration apply failed")
		return err
	}
	log.Info().Msg("migrations applied successfully")

	return nil
}
