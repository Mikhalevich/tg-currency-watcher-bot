package currencybot

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-telegram/bot"
	"github.com/google/uuid"

	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/domain/button"
	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/domain/user"
)

func (cb *CurrencyBot) UnsubscribeCurrencies(
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

	buttons, groupID, err := cb.makeUnsubscribeButtons(ctx, currencies)
	if err != nil {
		return fmt.Errorf("make buttons markup: %w", err)
	}

	if _, err := botAPI.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      info.ChatID,
		Text:        "choose pair to unsubscribe",
		ReplyMarkup: makeButtonGroupMarkup(groupID, buttons),
	}); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func (cb *CurrencyBot) makeUnsubscribeButtons(
	ctx context.Context,
	currencies []user.Currency,
) ([]button.Button, string, error) {
	var (
		groupID = uuid.NewString()
		buttons = make([]button.Button, 0, len(currencies))
	)

	for idx, currPair := range currencies {
		btn, err := button.UnsubscribeCurrencyPairButton(
			strconv.Itoa(idx+1),
			currPair.FormatPair(),
			button.UnsubscribeCurrencyPairPayload{
				CurrencyID:    currPair.ID,
				FormattedPair: currPair.FormatPair(),
			},
		)

		if err != nil {
			return nil, "", fmt.Errorf("make currency pair button: %w", err)
		}

		buttons = append(buttons, btn)
	}

	if err := cb.buttonProvider.SetButtonGroup(ctx, groupID, buttons); err != nil {
		return nil, "", fmt.Errorf("store buttons group: %w", err)
	}

	return buttons, groupID, nil
}
