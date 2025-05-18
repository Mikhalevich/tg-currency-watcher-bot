package user

import (
	"context"
	"fmt"
	"strings"
	"time"
)

func (u *UserProcessor) SendNotifications(ctx context.Context) error {
	if err := transaction(ctx, u.storage, func(ctx context.Context, store Storage) error {
		users, err := store.GetUsersReadyForNotifications(ctx)
		if err != nil {
			return fmt.Errorf("get users for notifications: %w", err)
		}

		userIDs := make([]int, 0, len(users))

		for _, user := range users {
			currencies, err := store.GetCurrenciesByChatID(ctx, user.ChatID)
			if err != nil {
				return fmt.Errorf("get user currencies: %w", err)
			}

			userIDs = append(userIDs, user.ID)

			u.sender.SendTextMessage(ctx, user.ChatID, formatUserCurrencies(currencies))
		}

		if err := store.UpdateLastNotificationTime(ctx, userIDs, time.Now()); err != nil {
			return fmt.Errorf("update last notification time: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}

func formatUserCurrencies(currencies []Currency) string {
	output := make([]string, 0, len(currencies))
	for _, v := range currencies {
		output = append(output, fmt.Sprintf("%s => %s", v.FormatPair(), v.FormatPrice()))
	}

	return strings.Join(output, "\n")
}
