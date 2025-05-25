package postgres

import (
	"context"
	"fmt"

	"github.com/volatiletech/sqlboiler/v4/queries"
)

func (p *Postgres) RemoveUserCurrency(ctx context.Context, userID int, currencyID int) error {
	res, err := queries.Raw(`
		DELETE FROM
			users_currency
		WHERE
			user_id = $1 AND
			currency_id = $2
	`, userID, currencyID).ExecContext(ctx, p.db)

	if err != nil {
		return fmt.Errorf("delete user currency: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return errNotFound
	}

	return nil
}
