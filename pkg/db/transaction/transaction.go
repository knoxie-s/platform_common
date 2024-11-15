package transaction

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/knoxie-s/platform_common/pkg/db"
	"github.com/knoxie-s/platform_common/pkg/db/pg"
	"github.com/pkg/errors"
)

type manager struct {
	db db.Transactor
}

// NewTxManager init
func NewTxManager(db db.Transactor) db.TxManager {
	return &manager{
		db: db,
	}
}

func (t *manager) transaction(ctx context.Context, options pgx.TxOptions, fn db.TxFnHandler) error {
	tx, ok := ctx.Value(pg.TxKey).(pgx.Tx)
	if ok {
		return fn(ctx)
	}

	tx, err := t.db.BeginTx(ctx, options)
	if err != nil {
		return err
	}

	ctx = pg.MakeContextTx(ctx, tx)

	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("panic recovered %v", r)
		}

		if err != nil {
			if errRollBack := tx.Rollback(ctx); errRollBack != nil {
				err = errors.Wrapf(err, "rollback error %v", errRollBack)
			}
		}

		if err = tx.Commit(ctx); err != nil {
			err = errors.Wrap(err, "failed commit transaction")
		}
	}()

	err = fn(ctx)
	if err != nil {
		return errors.Wrapf(err, "handle func err: %v", err.Error())
	}

	return nil
}

// ReadCommitted tx isolation
func (t *manager) ReadCommitted(ctx context.Context, fn db.TxFnHandler) error {
	trnLevel := pgx.TxOptions{IsoLevel: pgx.ReadCommitted}
	return t.transaction(ctx, trnLevel, fn)
}
