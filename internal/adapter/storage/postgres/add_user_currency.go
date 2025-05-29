package postgres

import (
	"context"
	"fmt"

	"github.com/volatiletech/sqlboiler/v4/queries"
)

func (p *Postgres) AddUserCurrency(ctx context.Context, userID int, currencyID int) error {
	res, err := queries.Raw(`
		INSERT INTO users_currency(
			user_id,
			currency_id
		) VALUES (
			$1,
			$2
		)
	`, userID, currencyID).ExecContext(ctx, p.db)

	if err != nil {
		if p.driver.IsConstraintError(err, "user_currency_pk") {
			return errAlreadyExists
		}

		return fmt.Errorf("insert users currency: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("no rows affected: %w", err)
	}

	return nil
}
