package postgres

import (
	"context"
	"fmt"

	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/domain/user"
)

func (p *Postgres) CreateUser(ctx context.Context, usr *user.User) (int, error) {
	dbUser := convertToDBUser(usr)

	if err := dbUser.Insert(ctx, p.db,
		boil.Whitelist("chat_id", "created_at", "notification_interval_hours", "last_notification_time"),
	); err != nil {
		return 0, fmt.Errorf("insert: %w", err)
	}

	return dbUser.ID, nil
}
