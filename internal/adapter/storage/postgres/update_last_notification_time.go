package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/volatiletech/sqlboiler/v4/queries"
)

func (p *Postgres) UpdateLastNotificationTime(
	ctx context.Context,
	userIDs []int,
	notificationTime time.Time,
) error {
	var (
		query = `
			UPDATE
				users
			SET
				last_notification_time = ?
			WHERE
				id IN (?)
		`
	)

	query, args, err := sqlx.In(query, notificationTime, userIDs)
	if err != nil {
		return fmt.Errorf("sqlx in: %w", err)
	}

	query = sqlx.Rebind(sqlx.DOLLAR, query)

	res, err := queries.Raw(query, args...).ExecContext(ctx, p.db)
	if err != nil {
		return fmt.Errorf("update users: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("no rows affected")
	}

	return nil
}
