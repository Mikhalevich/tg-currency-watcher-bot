package currencybot

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/domain/button"
	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/domain/rates"
)

func (cb *CurrencyBot) CurrencyPairs(ctx context.Context, botAPI *bot.Bot, update *models.Update) error {
	currencies, err := cb.ratesProvider.CurrencyRates(ctx)
	if err != nil {
		return fmt.Errorf("get user currencies: %w", err)
	}

	markup, err := cb.makeCurrencyPairsButtonsMarkup(ctx, currencies)
	if err != nil {
		return fmt.Errorf("make buttons markup: %w", err)
	}

	if _, err := botAPI.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        "choose pair to subscribe for notifications",
		ReplyMarkup: markup,
	}); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func (cb *CurrencyBot) makeCurrencyPairsButtonsMarkup(
	ctx context.Context,
	currencies []rates.Currency,
) (models.ReplyMarkup, error) {
	var (
		group, groupID = cb.buttonProvider.ButtonGroup()
		buttons        = make([]models.InlineKeyboardButton, 0, len(currencies))
	)

	for _, currPair := range currencies {
		btn, err := button.CurrencyPairButton(
			currPair.FormatPair(),
			button.CurrencyPairPayload{
				CurrencyID: currPair.ID,
				IsInverted: false,
			},
		)

		if err != nil {
			return nil, fmt.Errorf("make currency pair button: %w", err)
		}

		group.AddButton(btn)

		buttons = append(buttons, models.InlineKeyboardButton{
			Text:         btn.Caption,
			CallbackData: fmt.Sprintf("%s_%s", groupID, strconv.Itoa(currPair.ID)),
		})
	}

	if err := group.Store(ctx); err != nil {
		return nil, fmt.Errorf("store buttons group: %w", err)
	}

	return models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{buttons},
	}, nil
}
