package txmanager

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/ua-academy-projects/share-bite/pkg/database"
	"github.com/ua-academy-projects/share-bite/pkg/database/pg"
	"github.com/ua-academy-projects/share-bite/pkg/errwrap"
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
		return errwrap.Wrap("can't begin transaction", err)
	}

	ctx = pg.MakeContextTx(ctx, tx)

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic recovered: %v", r)
		}

		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				err = errwrap.Wrap("transaction rollback", rollbackErr)
			}

			return
		}

		if nil == err {
			if commitErr := tx.Commit(ctx); commitErr != nil {
				err = errwrap.Wrap("transaction commit", err)
			}
		}
	}()

	if err = fn(ctx); err != nil {
		err = errwrap.Wrap("execute code inside the transaction", err)
	}
	return err
}
