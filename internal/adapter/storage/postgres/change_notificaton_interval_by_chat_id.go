package postgres

import (
	"context"
	"fmt"

	"github.com/volatiletech/sqlboiler/v4/queries"
)

func (p *Postgres) ChangeNotificationIntervalByChatID(
	ctx context.Context,
	chatID int64,
	interval int,
) error {
	query := `
			UPDATE
				users
			SET
				notification_interval_hours = $1
			WHERE
				chat_id = $2
		`

	res, err := queries.Raw(query, interval, chatID).ExecContext(ctx, p.db)
	if err != nil {
		return fmt.Errorf("update interval: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errNotFound
	}

	return nil
}
