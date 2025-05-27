package user

import (
	"context"
	"fmt"
)

func (u *UserProcessor) UnsubscribeCurrency(ctx context.Context, chatID int64, currencyID int) error {
	usr, err := u.userByChatID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("get user by chat_id: %w", err)
	}

	if err := u.storage.RemoveUserCurrency(ctx, usr.ID, currencyID); err != nil {
		if u.storage.IsNotFoundError(err) {
			return ErrCurrencyNotFound
		}

		return fmt.Errorf("remove user currency: %w", err)
	}

	return nil
}
