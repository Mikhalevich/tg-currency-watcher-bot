package user

import (
	"context"
	"fmt"
)

func (u *UserProcessor) ChangeNotificationInterval(
	ctx context.Context,
	chatID int64,
	interval int,
) error {
	if err := u.storage.ChangeNotificationIntervalByChatID(ctx, chatID, interval); err != nil {
		if u.storage.IsNotFoundError(err) {
			return ErrChatNotFound
		}

		return fmt.Errorf("storage change interval: %w", err)
	}

	return nil
}
