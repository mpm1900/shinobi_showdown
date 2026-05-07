package db

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context) (*pgxpool.Pool, error) {
	logger := slog.Default()

	pgdb := getenvDefault("POSTGRES_DB", "shinobi")
	pguser := getenvDefault("POSTGRES_USER", "postgres")
	pgpass := getenvDefault("POSTGRES_PASSWORD", "postgres")
	pghost := getenvDefault("POSTGRES_HOST", "localhost")
	pgport := getenvDefault("POSTGRES_PORT", "5432")
	pgsslmode := getenvDefault("POSTGRES_SSLMODE", "disable")

	connStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		pguser, pgpass, pgdb, pghost, pgport, pgsslmode,
	)

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		logger.Error("Error connecting to database", "err", err)
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		logger.Error("Error pinging database", "err", err)
		pool.Close()
		return nil, err
	}

	return pool, nil
}

func getenvDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
