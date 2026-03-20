package pg

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ua-academy-projects/share-bite/pkg/database"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
)

type key string

const TxKey key = "tx"

type pg struct {
	pool *pgxpool.Pool
}

var _ database.DB = (*pg)(nil)

func NewDB(pool *pgxpool.Pool) *pg {
	return &pg{pool: pool}
}

func (p *pg) Ping(ctx context.Context) error {
	return p.pool.Ping(ctx)
}

func (p *pg) BeginTx(ctx context.Context, txOpts pgx.TxOptions) (pgx.Tx, error) {
	tx, err := p.pool.BeginTx(ctx, txOpts)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (p *pg) QueryContext(ctx context.Context, q database.Query, args ...any) (pgx.Rows, error) {
	logQuery(ctx, q)

	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.Query(ctx, q.Sql, args...)
	}

	return p.pool.Query(ctx, q.Sql, args...)
}

func (p *pg) QueryRowContext(ctx context.Context, q database.Query, args ...any) pgx.Row {
	logQuery(ctx, q)

	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.QueryRow(ctx, q.Sql, args...)
	}

	return p.pool.QueryRow(ctx, q.Sql, args...)
}

func (p *pg) ExecContext(ctx context.Context, q database.Query, args ...any) (pgconn.CommandTag, error) {
	logQuery(ctx, q)

	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.Exec(ctx, q.Sql, args...)
	}

	return p.pool.Exec(ctx, q.Sql, args...)
}

func MakeContextTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, TxKey, tx)
}

func (p *pg) Close() {
	p.pool.Close()
}

func (p *pg) Pool() *pgxpool.Pool {
	return p.pool
}

func logQuery(ctx context.Context, q database.Query) {
	logger.DebugKV(
		ctx,
		"log query",
		"sql", q.Name,
		"query", q.Sql,
	)
}
