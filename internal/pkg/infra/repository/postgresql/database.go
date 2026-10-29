package postgresql

import (
	"context"
	"log"

	"github.com/go-faster/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

var errNoRowsAffected = errors.New("no rows affected")

type Database struct {
	pool *pgxpool.Pool
}

func MustCreate(ctx context.Context, dsn string) *Database {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("can't create pgxpool: %v", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		log.Fatalf("can't ping pg: %v", err)
	}

	return &Database{
		pool: pool,
	}
}

// QueryRow is expected to return at most 1 row.
// Получает первую строку, а остальные отбрасывает
func (d *Database) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	return d.pool.QueryRow(ctx, query, args...)
}

// Query returns pgx.Rows
func (d *Database) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	return d.pool.Query(ctx, query, args...)
}

// Exec executes provided sql
// Возвращает pgconn.CommandTag, у которого есть RowsAffected
func (d *Database) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	return d.pool.Exec(ctx, query, args...)
}
