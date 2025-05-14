package currencybot

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/domain/user"
)

func (cb *CurrencyBot) MyCurrencies(ctx context.Context, botAPI *bot.Bot, update *models.Update) error {
	currencies, err := cb.userCurrency.GetCurrenciesByChatID(ctx, update.Message.Chat.ID)
	if err != nil {
		return fmt.Errorf("get user currencies: %w", err)
	}

	if len(currencies) == 0 {
		cb.replyTextMessage(ctx, update.Message.Chat.ID, update.Message.ID, "no subscribed currencies")

		return nil
	}

	cb.replyTextMessage(ctx, update.Message.Chat.ID, update.Message.ID, formatUserCurrencies(currencies))

	return nil
}

func formatUserCurrencies(currencies []user.Currency) string {
	output := make([]string, 0, len(currencies))
	for _, v := range currencies {
		output = append(output, fmt.Sprintf("%s => %s", v.FormatPair(), v.FormatPrice()))
	}

	return strings.Join(output, "\n")
}
