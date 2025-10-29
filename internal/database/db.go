package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

func Connect() error {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return fmt.Errorf("DATABASE_URL ")
	}

	cfg, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return fmt.Errorf("error trying to parse DATABASE_URL: %w", err)
	}

	dbpool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return fmt.Errorf("error connecting to PostgreSQL: %w", err)
	}

	if err := dbpool.Ping(context.Background()); err != nil {
		return fmt.Errorf("could not ping the database: %w", err)
	}

	pool = dbpool
	log.Println("Connected to PostgreSQL successfully")
	return nil
}

func Close() {
	if pool != nil {
		pool.Close()
	}
}

func GetPool() *pgxpool.Pool {
	return pool
}
