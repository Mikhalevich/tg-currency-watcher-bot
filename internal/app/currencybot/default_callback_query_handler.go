package currencybot

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/domain/button"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (cb *CurrencyBot) DefaultCallbackQueryHandler(
	ctx context.Context,
	botAPI *bot.Bot,
	update *models.Update,
) error {
	btn, err := cb.buttonProvider.GetButton(ctx, update.CallbackQuery.Data)
	if err != nil {
		return fmt.Errorf("get button: %w", err)
	}

	var (
		chatID    = update.CallbackQuery.Message.Message.Chat.ID
		messageID = update.CallbackQuery.Message.Message.ID
	)

	switch btn.Type {
	case button.CurrencyPair:
		if err := cb.processCurrencyPair(ctx, btn, chatID, messageID); err != nil {
			return fmt.Errorf("process currency pair: %w", err)
		}

	default:
		cb.replyTextMessage(
			ctx,
			chatID,
			messageID,
			"command is not supported",
		)
	}

	return nil
}

func (cb *CurrencyBot) processCurrencyPair(
	ctx context.Context,
	btn *button.Button,
	chatID int64,
	messageID int,
) error {
	payload, err := button.GetPayload[button.CurrencyPairPayload](*btn)
	if err != nil {
		return fmt.Errorf("get payload for currency pair: %w", err)
	}

	cb.replyTextMessage(ctx, chatID, messageID, fmt.Sprintf("get curren currency pair id: %d", payload.CurrencyID))

	return nil
}
