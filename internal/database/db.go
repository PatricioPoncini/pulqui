package database

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

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

	if err := runMigrations(pool); err != nil {
		return fmt.Errorf("error executing migrations: %w", err)
	}

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

func runMigrations(pool *pgxpool.Pool) error {
	migrationsDir := "internal/database/migrations"

	err := filepath.WalkDir(migrationsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(d.Name(), ".sql") {
			return nil
		}

		sqlBytes, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error reading %s: %w", path, err)
		}

		sqlContent := string(sqlBytes)
		_, execErr := pool.Exec(context.Background(), sqlContent)
		if execErr != nil {
			return fmt.Errorf("error executing %s: %w", path, execErr)
		}

		log.Printf("Migration executed: %s\n", d.Name())
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
