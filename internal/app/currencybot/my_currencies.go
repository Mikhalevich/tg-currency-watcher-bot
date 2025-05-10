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
	currencies, err := cb.userCurrency.GetUserCurrencies(ctx)
	if err != nil {
		return fmt.Errorf("get user currencies: %w", err)
	}

	cb.replyTextMessage(ctx, update.Message.Chat.ID, update.Message.ID, formatUserCurrencies(currencies))

	return nil
}

func formatUserCurrencies(currencies []user.Currency) string {
	output := make([]string, 0, len(currencies))
	for _, v := range currencies {
		output = append(output, fmt.Sprintf("%s/%s => %9.2f", v.Base, v.Quote, v.Price))
	}

	return strings.Join(output, "\n")
}
