package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/adapter/storage/postgres/internal/transaction"
	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/domain/exchange"
	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/domain/rates"
	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/domain/user"
)

var _ exchange.Storage = (*Postgres)(nil)
var _ user.Storage = (*Postgres)(nil)
var _ rates.RatesProvider = (*Postgres)(nil)

type Postgres struct {
	db boil.ContextExecutor
}

func New(db boil.ContextExecutor) *Postgres {
	return &Postgres{
		db: db,
	}
}

func (p *Postgres) IsNotFoundError(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

func (p *Postgres) Transaction(
	ctx context.Context,
	txFn func(ctx context.Context, store any) error,
) error {
	db, ok := p.db.(*sql.DB)
	if !ok {
		return errors.New("not sql.DB object")
	}

	if err := transaction.Transaction(ctx, db, func(ctx context.Context, tx *sql.Tx) error {
		if err := txFn(ctx, New(tx)); err != nil {
			return fmt.Errorf("tx func: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}

type Transactional interface {
	Transaction(
		ctx context.Context,
		txFn func(ctx context.Context, store any) error,
	) error
}

func Transaction[T any](
	ctx context.Context,
	transact Transactional,
	txFn func(ctx context.Context, store T) error,
) error {
	if err := transact.Transaction(ctx, func(ctx context.Context, store any) error {
		storeT, ok := store.(T)
		if !ok {
			return errors.New("invalid object")
		}

		if err := txFn(ctx, storeT); err != nil {
			return fmt.Errorf("tx fn: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}
