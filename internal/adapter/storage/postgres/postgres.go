package postgres

import (
	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/domain/exchange"
)

var _ exchange.Storage = (*Postgres)(nil)

type Postgres struct {
}

func New() *Postgres {
	return &Postgres{}
}
