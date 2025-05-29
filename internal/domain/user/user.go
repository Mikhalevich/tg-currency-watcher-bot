package user

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrChatNotFound          = errors.New("chat not found")
	ErrCurrencyNotFound      = errors.New("currency not found")
	ErrCurrencyAlreadyExists = errors.New("currency already exists")
)

type User struct {
	ID                        int
	ChatID                    int64
	CreatedAt                 time.Time
	NotificationIntervalHours int
	LastNotificationTime      time.Time
}

//nolint:interfacebloat
type Storage interface {
	GetCurrenciesByChatID(ctx context.Context, chatID int64) ([]Currency, error)
	CreateUser(ctx context.Context, usr *User) (int, error)
	GetUserByChatID(ctx context.Context, chatID int64) (*User, error)
	AddUserCurrency(ctx context.Context, userID int, currencyID int) error
	RemoveUserCurrency(ctx context.Context, userID int, currencyID int) error
	ChangeNotificationIntervalByChatID(ctx context.Context, chatID int64, interval int) error

	GetUsersReadyForNotifications(ctx context.Context, until time.Time, limit int) ([]*User, error)
	UpdateLastNotificationTime(ctx context.Context, userIDs []int, notificationTime time.Time) error
	Transaction(ctx context.Context, txFn func(ctx context.Context, store any) error) error

	IsNotFoundError(err error) bool
	IsAlreadyExistsError(err error) bool
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

func transaction(
	ctx context.Context,
	transact Storage,
	txFn func(ctx context.Context, store Storage) error,
) error {
	if err := transact.Transaction(ctx, func(ctx context.Context, store any) error {
		storeT, ok := store.(Storage)
		if !ok {
			return errors.New("invalid object")
		}

		if err := txFn(ctx, storeT); err != nil {
			return fmt.Errorf("tx fn: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}
