package postgres

import (
	"context"

	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/domain/exchange"
)

func (p *Postgres) GetCurrencies(ctx context.Context) ([]exchange.Currency, error) {
	return []exchange.Currency{
		{
			Base:            "USD",
			BaseExternalID:  "2781",
			Quote:           "BTC",
			QuoteExternalID: "1",
		},
	}, nil
}
