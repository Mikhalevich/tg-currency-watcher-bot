package currencybot

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/domain/user"
)

func (cb *CurrencyBot) MyCurrencies(ctx context.Context, botAPI *bot.Bot, update *models.Update) {
	currencies, err := cb.userCurrency.GetUserCurrencies(ctx)
	if err != nil {
		cb.logger.WithError(err).Error("get user currencies")

		return
	}

	if _, err := botAPI.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		ReplyParameters: &models.ReplyParameters{
			MessageID: update.Message.ID,
		},
		Text: formatUserCurrencies(currencies),
	}); err != nil {
		cb.logger.WithError(err).Error("send message")
	}
}

func formatUserCurrencies(currencies []user.Currency) string {
	output := make([]string, 0, len(currencies))
	for _, v := range currencies {
		output = append(output, fmt.Sprintf("%s-%s => %9.2f", v.Quote, v.Base, 1/v.Price))
	}

	return strings.Join(output, "n")
}
