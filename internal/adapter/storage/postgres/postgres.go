package postgres

import (
	"database/sql"

	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/domain/exchange"
)

var _ exchange.Storage = (*Postgres)(nil)

type Postgres struct {
	db *sql.DB
}

func New(db *sql.DB) *Postgres {
	return &Postgres{
		db: db,
	}
}
