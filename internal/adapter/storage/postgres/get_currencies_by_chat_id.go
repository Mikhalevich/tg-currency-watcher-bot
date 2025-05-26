package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/adapter/storage/postgres/internal/models"
	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/domain/user"
)

func (p *Postgres) GetCurrenciesByChatID(ctx context.Context, chatID int64) ([]user.Currency, error) {
	dbUsr, err := models.Users(
		qm.Load(models.UserRels.Currencies),
		qm.Where("chat_id = ?", chatID),
	).One(ctx, p.db)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("get user currencies: %w", err)
	}

	return toUserCurrencies(dbUsr.R.Currencies), nil
}

func toUserCurrencies(dbCurrencies []*models.Currency) []user.Currency {
	if len(dbCurrencies) == 0 {
		return nil
	}

	currencies := make([]user.Currency, 0, len(dbCurrencies))

	for _, v := range dbCurrencies {
		currencies = append(currencies, toUserCurrency(v))
	}

	return currencies
}

func toUserCurrency(dbCurrency *models.Currency) user.Currency {
	price, _ := dbCurrency.Price.Float64()

	return user.Currency{
		ID:         dbCurrency.ID,
		Base:       user.Symbol(dbCurrency.Base),
		Quote:      user.Symbol(dbCurrency.Quote),
		Price:      user.Money(price),
		IsInverted: dbCurrency.IsInverted,
		UpdatedAt:  dbCurrency.UpdatedAt,
	}
}
