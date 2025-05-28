package currencybot

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/google/uuid"

	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/domain/button"
	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/domain/rates"
)

const (
	buttonsPerRow = 4
)

func (cb *CurrencyBot) CurrencyPairs(
	ctx context.Context,
	botAPI *bot.Bot,
	info MessageInfo,
) error {
	currencies, err := cb.ratesProvider.CurrencyRates(ctx)
	if err != nil {
		return fmt.Errorf("get user currencies: %w", err)
	}

	if len(currencies) == 0 {
		cb.replyTextMessage(ctx, info.ChatID, info.MessageID, "no currency pairs")

		return nil
	}

	buttons, groupID, err := cb.makeCurrencyPairsButtons(ctx, currencies)
	if err != nil {
		return fmt.Errorf("make buttons markup: %w", err)
	}

	if _, err := botAPI.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      info.ChatID,
		Text:        "choose pair to subscribe for notifications",
		ReplyMarkup: makeButtonGroupMarkup(groupID, buttons),
	}); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func (cb *CurrencyBot) makeCurrencyPairsButtons(
	ctx context.Context,
	currencies []rates.Currency,
) ([]button.Button, string, error) {
	var (
		groupID = uuid.NewString()
		buttons = make([]button.Button, 0, len(currencies))
	)

	for idx, currPair := range currencies {
		btn, err := button.CurrencyPairButton(
			strconv.Itoa(idx+1),
			currPair.FormatPair(),
			button.CurrencyPairPayload{
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

func makeButtonGroupMarkup(groupID string, buttons []button.Button) models.ReplyMarkup {
	var (
		buttonRow             = make([]models.InlineKeyboardButton, 0, len(buttons))
		rows                  = len(buttons) / buttonsPerRow
		markup                = make([][]models.InlineKeyboardButton, 0, rows)
		prevButtonRowStartIdx = 0
	)

	for idx, btn := range buttons {
		buttonRow = append(buttonRow, models.InlineKeyboardButton{
			Text:         btn.Caption,
			CallbackData: fmt.Sprintf("%s_%s", groupID, btn.ID),
		})

		if (idx+1)%buttonsPerRow == 0 {
			markup = append(markup, buttonRow[prevButtonRowStartIdx:])
			prevButtonRowStartIdx = idx + 1
		}
	}

	if prevButtonRowStartIdx < len(buttonRow) {
		markup = append(markup, buttonRow[prevButtonRowStartIdx:])
	}

	return models.InlineKeyboardMarkup{
		InlineKeyboard: markup,
	}
}
