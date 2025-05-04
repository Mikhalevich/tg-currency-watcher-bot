package postgres

import (
	"context"

	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/domain/exchange"
)

func (p *Postgres) UpdateCurrencies(ctx context.Context, currencies []exchange.Currency) error {
	return nil
}
