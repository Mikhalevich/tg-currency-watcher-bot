package postgres

import (
	"context"
	"fmt"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/adapter/storage/postgres/internal/models"
	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/domain/user"
)

func (p *Postgres) GetUserByChatID(ctx context.Context, chatID int64) (*user.User, error) {
	dbUser, err := models.Users(qm.Where("chat_id = ?", chatID)).One(ctx, p.db)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	return convertToUser(dbUser), nil
}

func convertToUser(dbUser *models.User) *user.User {
	return &user.User{
		ID:                        dbUser.ID,
		ChatID:                    int64(dbUser.ChatID),
		CreatedAt:                 dbUser.CreatedAt,
		NotificationIntervalHours: dbUser.NotificationIntervalHours,
		LastNotificationTime:      dbUser.LastNotificationTime,
	}
}

func convertToDBUser(usr *user.User) *models.User {
	return &models.User{
		ID:                        usr.ID,
		ChatID:                    int(usr.ChatID),
		CreatedAt:                 usr.CreatedAt,
		NotificationIntervalHours: usr.NotificationIntervalHours,
		LastNotificationTime:      usr.LastNotificationTime,
	}
}
