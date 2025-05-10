package postgres

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/adapter/storage/postgres/internal/models"
	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/domain/exchange"
)

func (p *Postgres) GetCurrencies(ctx context.Context) ([]exchange.Currency, error) {
	currencies, err := models.Currencies().All(ctx, p.db)
	if err != nil {
		return nil, fmt.Errorf("get all currencies: %w", err)
	}

	return toPortCurrencies(currencies), nil
}

func toPortCurrencies(dbCurrencies []*models.Currency) []exchange.Currency {
	currencies := make([]exchange.Currency, 0, len(dbCurrencies))

	for _, v := range dbCurrencies {
		currencies = append(currencies, toPortCurrency(v))
	}

	return currencies
}

func toPortCurrency(dbCurrency *models.Currency) exchange.Currency {
	price, _ := dbCurrency.Price.Float64()

	return exchange.Currency{
		ID:              dbCurrency.ID,
		Base:            exchange.Symbol(dbCurrency.Base),
		BaseExternalID:  exchange.ExternalID(dbCurrency.BaseExternalID),
		Quote:           exchange.Symbol(dbCurrency.Quote),
		QuoteExternalID: exchange.ExternalID(dbCurrency.QuoteExternalID),
		Price:           exchange.Money(price),
		IsInverted:      dbCurrency.IsInverted,
		UpdatedAt:       dbCurrency.UpdatedAt,
	}
}
