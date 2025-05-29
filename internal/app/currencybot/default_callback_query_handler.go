package currencybot

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-telegram/bot"

	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/domain/button"
	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/domain/user"
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
		if err := cb.processCurrencyPair(ctx, btn, info.ChatID); err != nil {
			return fmt.Errorf("process currency pair: %w", err)
		}

	case button.UnsubscribeCurrencyPair:
		if err := cb.processUnsubscribeCurrencyPair(ctx, btn, info.ChatID); err != nil {
			return fmt.Errorf("process unsubscribe currency pair: %w", err)
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
) error {
	payload, err := button.GetPayload[button.CurrencyPairPayload](*btn)
	if err != nil {
		return fmt.Errorf("get payload for currency pair: %w", err)
	}

	if err := cb.userCurrency.SubscribeCurrency(ctx, chatID, payload.CurrencyID); err != nil {
		if errors.Is(err, user.ErrCurrencyAlreadyExists) {
			cb.sendTextMessage(ctx, chatID, fmt.Sprintf("Currency %s already subscribed", payload.FormattedPair))

			return nil
		}

		return fmt.Errorf("subscribe currency: %w", err)
	}

	cb.sendTextMessage(ctx, chatID, fmt.Sprintf("Subscibed %s", payload.FormattedPair))

	return nil
}

func (cb *CurrencyBot) processUnsubscribeCurrencyPair(
	ctx context.Context,
	btn *button.Button,
	chatID int64,
) error {
	payload, err := button.GetPayload[button.UnsubscribeCurrencyPairPayload](*btn)
	if err != nil {
		return fmt.Errorf("get payload for unsubscribe currency pair: %w", err)
	}

	if err := cb.userCurrency.UnsubscribeCurrency(ctx, chatID, payload.CurrencyID); err != nil {
		if errors.Is(err, user.ErrCurrencyNotFound) {
			cb.sendTextMessage(ctx, chatID, "No such currency pair")

			return nil
		}

		return fmt.Errorf("unsubscribe currency: %w", err)
	}

	cb.sendTextMessage(ctx, chatID, fmt.Sprintf("Unsubscibed %s", payload.FormattedPair))

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
		if errors.Is(err, user.ErrChatNotFound) {
			cb.sendTextMessage(ctx, chatID, "no subscribed currencies")

			return nil
		}

		return fmt.Errorf("change interval: %w", err)
	}

	cb.sendTextMessage(ctx, chatID, "interval changed")

	return nil
}
