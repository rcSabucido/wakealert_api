package db

import (
	"context"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rcsabucido/wakealert_api/internal/db/sqlc"
)

var (
	pool *pgxpool.Pool
	once sync.Once
)

func Queries(ctx context.Context) (*sqlc.Queries, error) {
	var err error
	once.Do(func() {
		pool, err = pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	})
	if err != nil {
		return nil, err
	}
	return sqlc.New(pool), nil
}

func Pool(ctx context.Context) (*pgxpool.Pool, error) {
	var err error
	once.Do(func() {
		pool, err = pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	})
	if err != nil {
		return nil, err
	}
	return pool, nil
}
