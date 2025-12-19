package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Init(ctx context.Context) error {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return fmt.Errorf("DATABASE_URL is not set; set it in environment or .env")
	}

	cfg, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return err
	}
	cfg.MaxConns = 5
	cfg.MaxConnLifetime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return err
	}

	Pool = pool

	return nil
}

func Close() {
	if Pool != nil {
		Pool.Close()
	}
}
