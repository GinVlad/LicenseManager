package database

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(databaseURL string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse database URL: %w", err)
	}

	cfg.MaxConns = 20
	cfg.MinConns = 2
	cfg.MaxConnLifetime = time.Hour
	cfg.MaxConnIdleTime = 30 * time.Minute

	var pool *pgxpool.Pool
	var pingErr error

	// Retry connecting up to 5 times (useful for Railway/Docker when DB is booting)
	for i := 0; i < 5; i++ {
		pool, err = pgxpool.NewWithConfig(context.Background(), cfg)
		if err == nil {
			pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
			pingErr = pool.Ping(pingCtx)
			pingCancel()
			if pingErr == nil {
				break
			}
			pool.Close()
		}
		fmt.Printf("Database not ready yet, retrying in 3 seconds... (%d/5)\n", i+1)
		time.Sleep(3 * time.Second)
	}

	if pingErr != nil {
		return nil, fmt.Errorf("ping database failed after retries: %w", pingErr)
	}

	return pool, nil
}

func RunMigrations(pool *pgxpool.Pool) error {
	ctx := context.Background()

	// Create migration tracking table
	_, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMPTZ DEFAULT NOW()
		)`)
	if err != nil {
		return fmt.Errorf("create schema_migrations: %w", err)
	}

	// Collect applied versions
	rows, err := pool.Query(ctx, "SELECT version FROM schema_migrations ORDER BY version")
	if err != nil {
		return fmt.Errorf("query applied migrations: %w", err)
	}
	applied := map[string]bool{}
	for rows.Next() {
		var v string
		rows.Scan(&v)
		applied[v] = true
	}
	rows.Close()

	// Find all up migrations
	matches, err := filepath.Glob("internal/database/migrations/*.up.sql")
	if err != nil {
		return fmt.Errorf("glob migrations: %w", err)
	}
	sort.Strings(matches)

	for _, path := range matches {
		version := strings.TrimSuffix(filepath.Base(path), ".up.sql")
		if applied[version] {
			continue
		}

		sql, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}

		if _, err := pool.Exec(ctx, string(sql)); err != nil {
			return fmt.Errorf("apply %s: %w", version, err)
		}

		pool.Exec(ctx, "INSERT INTO schema_migrations (version) VALUES ($1)", version)
		fmt.Printf("migration applied: %s\n", version)
	}

	return nil
}
