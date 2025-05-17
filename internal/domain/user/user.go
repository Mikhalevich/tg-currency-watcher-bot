package user

import (
	"context"
	"time"
)

type User struct {
	ID                        int
	ChatID                    int64
	CreatedAt                 time.Time
	NotificationIntervalHours int
	LastNotificationTime      time.Time
}

type Storage interface {
	GetCurrenciesByChatID(ctx context.Context, chatID int64) ([]Currency, error)
	CreateUser(ctx context.Context, usr *User) (int, error)
	GetUserByChatID(ctx context.Context, chatID int64) (*User, error)
	AddUserCurrency(ctx context.Context, userID int, currencyID int) error

	GetUsersReadyForNotifications(ctx context.Context) ([]User, error)
	UpdateLastNotificationTime(ctx context.Context, userIDs []int, notificationTime time.Time) error

	IsNotFoundError(err error) bool
}

type MessageSender interface {
	SendTextMessage(ctx context.Context, chatID int64, text string)
}

type UserProcessor struct {
	storage Storage
	sender  MessageSender
}

func NewProcessor(storage Storage, sender MessageSender) *UserProcessor {
	return &UserProcessor{
		storage: storage,
		sender:  sender,
	}
}
