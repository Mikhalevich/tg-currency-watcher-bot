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

	IsNotFoundError(err error) bool
}

type UserProcessor struct {
	storage Storage
}

func NewProcessor(storage Storage) *UserProcessor {
	return &UserProcessor{
		storage: storage,
	}
}
