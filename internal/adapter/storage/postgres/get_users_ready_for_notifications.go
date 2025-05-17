package postgres

import (
	"context"

	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/domain/user"
)

func (p *Postgres) GetUsersReadyForNotifications(ctx context.Context) ([]user.User, error) {
	return nil, nil
}
