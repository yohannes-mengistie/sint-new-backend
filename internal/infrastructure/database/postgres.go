package database

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yohannes/sint-new-backend/internal/application"
)

type Postgres struct {
	pool *pgxpool.Pool
}

var _ interface {
	application.DatabaseHealth
	application.UserRepository
	application.ProductRepository
	application.PriceUpdateRepository
} = (*Postgres)(nil)

func NewPostgres(ctx context.Context, connectionString string) (*Postgres, error) {
	pool, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return &Postgres{pool: pool}, nil
}

func (database *Postgres) Ping(ctx context.Context) error {
	return database.pool.Ping(ctx)
}

func (database *Postgres) Close() {
	database.pool.Close()
}

func (database *Postgres) ApplyMigrations(ctx context.Context, directory string) error {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return fmt.Errorf("read migrations directory: %w", err)
	}
	files := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() &&
			strings.HasSuffix(entry.Name(), ".sql") &&
			!strings.Contains(entry.Name(), "_seed_") {
			files = append(files, filepath.Join(directory, entry.Name()))
		}
	}
	sort.Strings(files)
	for _, file := range files {
		sql, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", file, err)
		}
		if _, err := database.pool.Exec(ctx, string(sql)); err != nil {
			return fmt.Errorf("apply migration %s: %w", file, err)
		}
	}
	return nil
}
