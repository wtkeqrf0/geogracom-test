package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(ctx context.Context, conn string) *pgxpool.Pool {
	pool, err := pgxpool.New(ctx, conn)
	if err != nil {
		panic(err)
	}

	if err = pool.Ping(ctx); err != nil {
		panic(err)
	}

	context.AfterFunc(ctx, pool.Close)
	return pool
}
