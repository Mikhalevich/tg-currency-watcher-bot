package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/volatiletech/sqlboiler/v4/queries"

	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/adapter/storage/postgres/internal/models"
	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/domain/user"
)

func (p *Postgres) GetUsersReadyForNotifications(
	ctx context.Context,
	until time.Time,
	limit int,
) ([]*user.User, error) {
	var (
		query = `
			SELECT
				id,
				chat_id,
				created_at,
				notification_interval_hours,
				last_notification_time
			FROM
				users
			WHERE
				next_notification_time < $1::timestamptz
			LIMIT
				$2
			FOR UPDATE SKIP LOCKED
		`

		users models.UserSlice
	)

	if err := queries.Raw(query, until, limit).Bind(ctx, p.db, &users); err != nil {
		return nil, fmt.Errorf("select users: %w", err)
	}

	return convertToUsers(users), nil
}
