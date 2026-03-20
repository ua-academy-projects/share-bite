package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Client interface {
	DB() DB
	Close()
}

type DB interface {
	Pinger
	Transactor
	SQLExecer
	Close()
	Pool() *pgxpool.Pool
}

type Pinger interface {
	Ping(ctx context.Context) error
}

type Transactor interface {
	BeginTx(ctx context.Context, txOpts pgx.TxOptions) (pgx.Tx, error)
}

type SQLExecer interface {
	QueryExecer
}

type QueryExecer interface {
	QueryContext(ctx context.Context, q Query, args ...any) (pgx.Rows, error)
	QueryRowContext(ctx context.Context, q Query, args ...any) pgx.Row
	ExecContext(ctx context.Context, q Query, args ...any) (pgconn.CommandTag, error)
}

type Query struct {
	Name string
	Sql  string
}

type TxManager interface {
	ReadCommited(ctx context.Context, fn Handler) error
}

type Handler func(ctx context.Context) error
