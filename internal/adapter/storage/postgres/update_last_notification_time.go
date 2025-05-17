package postgres

import (
	"context"
	"time"
)

func (p *Postgres) UpdateLastNotificationTime(
	ctx context.Context,
	userIDs []int,
	notificationTime time.Time,
) error {
	return nil
}
