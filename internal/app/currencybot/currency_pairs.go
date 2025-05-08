package currencybot

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/domain/rates"
)

func (cb *CurrencyBot) CurrencyPairs(ctx context.Context, botAPI *bot.Bot, update *models.Update) error {
	currencies, err := cb.ratesProvider.CurrencyRates(ctx)
	if err != nil {
		return fmt.Errorf("get user currencies: %w", err)
	}

	cb.replyTextMessage(ctx, update.Message.Chat.ID, update.Message.ID, formatCurrencyPairs(currencies))

	return nil
}

func formatCurrencyPairs(currencies []rates.Currency) string {
	pairs := make([]string, 0, len(currencies))
	for _, v := range currencies {
		pairs = append(pairs, fmt.Sprintf("%s/%s", v.Base, v.Quote))
	}

	return strings.Join(pairs, "\n")
}
