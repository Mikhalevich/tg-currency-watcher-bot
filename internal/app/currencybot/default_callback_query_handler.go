package currencybot

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"

	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/domain/button"
)

func (cb *CurrencyBot) DefaultCallbackQueryHandler(
	ctx context.Context,
	botAPI *bot.Bot,
	info MessageInfo,
) error {
	btn, err := cb.buttonProvider.GetButton(ctx, info.Data)
	if err != nil {
		return fmt.Errorf("get button: %w", err)
	}

	switch btn.Type {
	case button.CurrencyPair:
		if err := cb.processCurrencyPair(ctx, btn, info.ChatID, info.MessageID); err != nil {
			return fmt.Errorf("process currency pair: %w", err)
		}

	case button.NotificationInterval:
		if err := cb.processNotificationInterval(ctx, btn, info.ChatID); err != nil {
			return fmt.Errorf("process notification interval: %w", err)
		}

	default:
		cb.replyTextMessage(
			ctx,
			info.ChatID,
			info.MessageID,
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

	if err := cb.userCurrency.SubscribeCurrency(ctx, chatID, payload.CurrencyID); err != nil {
		return fmt.Errorf("subscribe currency: %w", err)
	}

	cb.replyTextMessage(ctx, chatID, messageID, "Subscibed successfully")

	return nil
}

func (cb *CurrencyBot) processNotificationInterval(
	ctx context.Context,
	btn *button.Button,
	chatID int64,
) error {
	payload, err := button.GetPayload[button.NotificationIntervalPayload](*btn)
	if err != nil {
		return fmt.Errorf("get payload: %w", err)
	}

	if err := cb.userCurrency.ChangeNotificationInterval(ctx, chatID, payload.Interval); err != nil {
		return fmt.Errorf("change interval: %w", err)
	}

	return nil
}
