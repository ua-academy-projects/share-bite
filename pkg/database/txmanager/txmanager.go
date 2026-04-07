package txmanager

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/ua-academy-projects/share-bite/pkg/database"
	"github.com/ua-academy-projects/share-bite/pkg/database/pg"
)

type manager struct {
	db database.Transactor
}

func NewTransactionManager(db database.Transactor) database.TxManager {
	return &manager{
		db: db,
	}
}

func (m *manager) ReadCommited(ctx context.Context, fn database.Handler) error {
	return m.transaction(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}, fn)
}

func (m *manager) transaction(ctx context.Context, opts pgx.TxOptions, fn database.Handler) (err error) {
	tx, ok := ctx.Value(pg.TxKey).(pgx.Tx)
	if ok {
		return fn(ctx)
	}

	tx, err = m.db.BeginTx(ctx, opts)
	if err != nil {
		return fmt.Errorf("can't begin transaction: %w", err)
	}

	ctx = pg.MakeContextTx(ctx, tx)

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic recovered: %v", r)
		}

		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				err = fmt.Errorf("transaction rollback: %w", rollbackErr)
			}

			return
		}

		if nil == err {
			if commitErr := tx.Commit(ctx); commitErr != nil {
				err = fmt.Errorf("transaction commit: %w", err)
			}
		}
	}()

	if err = fn(ctx); err != nil {
		err = fmt.Errorf("execute code inside the transaction: %w", err)
	}
	return err
}
