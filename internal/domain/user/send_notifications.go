package user

import (
	"context"
	"fmt"
	"strings"
	"time"
)

func (u *UserProcessor) SendNotifications(ctx context.Context, usersLimit int) error {
	if err := transaction(ctx, u.storage, func(ctx context.Context, store Storage) error {
		currTime := time.Now()

		users, err := store.GetUsersReadyForNotifications(ctx, currTime, usersLimit)
		if err != nil {
			return fmt.Errorf("get users for notifications: %w", err)
		}

		if len(users) == 0 {
			return nil
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

		if err := store.UpdateLastNotificationTime(ctx, userIDs, currTime); err != nil {
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
