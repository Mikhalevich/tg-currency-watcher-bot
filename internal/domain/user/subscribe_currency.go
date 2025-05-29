package user

import (
	"context"
	"fmt"
	"time"
)

func (u *UserProcessor) SubscribeCurrency(ctx context.Context, chatID int64, currencyID int) error {
	usr, err := u.userByChatID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("get user by chat_id: %w", err)
	}

	if err := u.storage.AddUserCurrency(ctx, usr.ID, currencyID); err != nil {
		if u.storage.IsAlreadyExistsError(err) {
			return ErrCurrencyAlreadyExists
		}

		return fmt.Errorf("add user currency: %w", err)
	}

	return nil
}

func (u *UserProcessor) userByChatID(ctx context.Context, chatID int64) (*User, error) {
	usr, err := u.storage.GetUserByChatID(ctx, chatID)
	if err != nil {
		if u.storage.IsNotFoundError(err) {
			usr, err := u.createUser(ctx, chatID)
			if err != nil {
				return nil, fmt.Errorf("create user: %w", err)
			}

			return usr, nil
		}

		return nil, fmt.Errorf("get user by chat_id: %w", err)
	}

	return usr, nil
}

func (u *UserProcessor) createUser(ctx context.Context, chatID int64) (*User, error) {
	var (
		now = time.Now()
		usr = &User{
			ChatID:                    chatID,
			CreatedAt:                 now,
			NotificationIntervalHours: 1,
			LastNotificationTime:      now,
		}
	)

	userID, err := u.storage.CreateUser(ctx, usr)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	usr.ID = userID

	return usr, nil
}
