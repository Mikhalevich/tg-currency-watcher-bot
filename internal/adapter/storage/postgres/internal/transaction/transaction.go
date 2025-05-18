package transaction

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type TxFunc func(ctx context.Context, tx *sql.Tx) error

func Transaction(
	ctx context.Context,
	db *sql.DB,
	txFn TxFunc,
) error {
	trx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	//nolint:errcheck
	defer trx.Rollback()

	if err := txFn(ctx, trx); err != nil {
		if rollbackErr := trx.Rollback(); rollbackErr != nil {
			return errors.Join(fmt.Errorf("tx body: %w", err), fmt.Errorf("rollback: %w", err))
		}

		return fmt.Errorf("tx body: %w", err)
	}

	if err := trx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}
