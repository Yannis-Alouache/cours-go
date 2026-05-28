package db

import (
	"context"
	"fmt"
	"io/fs"
	"sort"

	sqldata "cours-go/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("create pgx pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return pool, nil
}

func RunMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	entries, err := fs.ReadDir(sqldata.FS, "migrations")
	if err != nil {
		return fmt.Errorf("read migrations: %w", err)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		content, err := sqldata.FS.ReadFile("migrations/" + entry.Name())
		if err != nil {
			return fmt.Errorf("read migration %s: %w", entry.Name(), err)
		}

		if _, err := pool.Exec(ctx, string(content)); err != nil {
			return fmt.Errorf("exec migration %s: %w", entry.Name(), err)
		}
	}

	return nil
}

func Seed(ctx context.Context, pool *pgxpool.Pool) error {
	content, err := sqldata.FS.ReadFile("seed.sql")
	if err != nil {
		return fmt.Errorf("read seed: %w", err)
	}

	if _, err := pool.Exec(ctx, string(content)); err != nil {
		return fmt.Errorf("exec seed: %w", err)
	}

	return nil
}
