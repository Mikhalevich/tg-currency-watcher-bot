package currencybot

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"

	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/domain/user"
)

func (cb *CurrencyBot) MyCurrencies(
	ctx context.Context,
	botAPI *bot.Bot,
	info MessageInfo,
) error {
	currencies, err := cb.userCurrency.GetCurrenciesByChatID(ctx, info.ChatID)
	if err != nil {
		return fmt.Errorf("get user currencies: %w", err)
	}

	if len(currencies) == 0 {
		cb.replyTextMessage(ctx, info.ChatID, info.MessageID, "no subscribed currencies")

		return nil
	}

	cb.replyTextMessage(ctx, info.ChatID, info.MessageID, formatUserCurrencies(currencies))

	return nil
}

func formatUserCurrencies(currencies []user.Currency) string {
	output := make([]string, 0, len(currencies))
	for _, v := range currencies {
		output = append(output, fmt.Sprintf("%s => %s", v.FormatPair(), v.FormatPrice()))
	}

	return strings.Join(output, "\n")
}
