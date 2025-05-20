package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/volatiletech/sqlboiler/v4/queries"
)

func (p *Postgres) UpdateLastNotificationTime(
	ctx context.Context,
	userIDs []int,
	notificationTime time.Time,
) error {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	query, args, err := psql.Update("users").
		Set("last_notification_time", notificationTime).
		Where(squirrel.Eq{
			"id": userIDs,
		}).ToSql()

	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

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
