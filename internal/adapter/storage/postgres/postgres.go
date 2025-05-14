package postgres

import (
	"database/sql"
	"errors"

	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/domain/exchange"
	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/domain/rates"
	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/domain/user"
)

var _ exchange.Storage = (*Postgres)(nil)
var _ user.Storage = (*Postgres)(nil)
var _ rates.RatesProvider = (*Postgres)(nil)

type Postgres struct {
	db *sql.DB
}

func New(db *sql.DB) *Postgres {
	return &Postgres{
		db: db,
	}
}

func (p *Postgres) IsNotFoundError(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
