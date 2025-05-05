package postgres

import (
	"context"
	"fmt"

	"github.com/ericlagergren/decimal"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/types"

	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/adapter/storage/postgres/internal/models"
	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/domain/exchange"
)

func (p *Postgres) UpdateCurrencies(ctx context.Context, currencies []exchange.Currency) error {
	dbCurrencies := toDBCurrencies(currencies)

	for _, v := range dbCurrencies {
		if _, err := v.Update(ctx, p.db, boil.Infer()); err != nil {
			return fmt.Errorf("update row: %w", err)
		}
	}

	return nil
}

func toDBCurrencies(currencies []exchange.Currency) []*models.Currency {
	dbCurrencies := make([]*models.Currency, 0, len(currencies))

	for _, v := range currencies {
		dbCurrencies = append(dbCurrencies, toDBCurrency(v))
	}

	return dbCurrencies
}

func toDBCurrency(currency exchange.Currency) *models.Currency {
	price := types.NewDecimal(decimal.New(0, 0).SetFloat64(float64(currency.Price)))

	return &models.Currency{
		ID:              currency.ID,
		Base:            currency.Base.String(),
		BaseExternalID:  currency.BaseExternalID.String(),
		Quote:           currency.Quote.String(),
		QuoteExternalID: currency.QuoteExternalID.String(),
		Price:           price,
		UpdatedAt:       currency.UpdatedAt,
	}
}
