package postgres

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/adapter/storage/postgres/internal/models"
	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/domain/rates"
)

func (p *Postgres) GetCurrencyRates(ctx context.Context) ([]rates.Currency, error) {
	currencies, err := models.Currencies().All(ctx, p.db)
	if err != nil {
		return nil, fmt.Errorf("get all currencies: %w", err)
	}

	return toRatesCurrencies(currencies), nil
}

func toRatesCurrencies(dbCurrencies []*models.Currency) []rates.Currency {
	currencies := make([]rates.Currency, 0, len(dbCurrencies))

	for _, v := range dbCurrencies {
		currencies = append(currencies, toRatesCurrency(v))
	}

	return currencies
}

func toRatesCurrency(dbCurrency *models.Currency) rates.Currency {
	price, _ := dbCurrency.Price.Float64()

	return rates.Currency{
		ID:        dbCurrency.ID,
		Base:      rates.Symbol(dbCurrency.Base),
		Quote:     rates.Symbol(dbCurrency.Quote),
		Price:     rates.Money(price),
		UpdatedAt: dbCurrency.UpdatedAt,
	}
}
