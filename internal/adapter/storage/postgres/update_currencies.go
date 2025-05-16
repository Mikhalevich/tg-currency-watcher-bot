package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/domain/exchange"
)

const (
	insertCurrencyColumns = 7
)

func (p *Postgres) UpdateCurrencies(ctx context.Context, currencies []exchange.Currency) error {
	var (
		queryTemplate = `
			INSERT INTO currency(
				base,
				base_external_id,
				quote,
				quote_external_id,
				price,
				is_inverted,
				updated_at
			) VALUES %s ON CONFLICT(base_external_id, quote_external_id) DO UPDATE SET 
				price = EXCLUDED.price,
				updated_at = EXCLUDED.updated_at
		`

		placeholders = make([]string, 0, len(currencies))
		values       = make([]any, 0, len(currencies)*insertCurrencyColumns)
	)

	for i, curr := range currencies {
		index := i * insertCurrencyColumns
		placeholders = append(placeholders,
			fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d)",
				//nolint:mnd
				index+1, index+2, index+3, index+4, index+5, index+6, index+7),
		)

		values = append(values,
			curr.Base,
			curr.BaseExternalID,
			curr.Quote,
			curr.QuoteExternalID,
			curr.Price,
			curr.IsInverted,
			curr.UpdatedAt,
		)
	}

	query := fmt.Sprintf(queryTemplate, strings.Join(placeholders, ","))

	res, err := p.db.ExecContext(ctx, query, values...)
	if err != nil {
		return fmt.Errorf("exec context: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return errors.New("no rows affected")
	}

	return nil
}
