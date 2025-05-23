package currencybot

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-telegram/bot"
	"github.com/google/uuid"

	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/domain/button"
)

func (cb *CurrencyBot) NotificationInterval(
	ctx context.Context,
	botAPI *bot.Bot,
	info MessageInfo,
) error {
	buttons, groupID, err := cb.makeNotificationIntervalButtons(
		ctx,
		[]int{1, 2, 3, 4, 5, 6, 7, 8, 12, 16, 20, 24},
	)

	if err != nil {
		return fmt.Errorf("make notification interval buttons: %w", err)
	}

	if _, err := botAPI.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      info.ChatID,
		Text:        "change notification interval",
		ReplyMarkup: makeButtonGroupMarkup(groupID, buttons),
	}); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func (cb *CurrencyBot) makeNotificationIntervalButtons(
	ctx context.Context,
	intervals []int,
) ([]button.Button, string, error) {
	var (
		groupID = uuid.NewString()
		buttons = make([]button.Button, 0, len(intervals))
	)

	for idx, interval := range intervals {
		btn, err := button.NotificationIntervalButton(
			strconv.Itoa(idx+1),
			strconv.Itoa(interval),
			button.NotificationIntervalPayload{
				Interval: interval,
			},
		)

		if err != nil {
			return nil, "", fmt.Errorf("make button: %w", err)
		}

		buttons = append(buttons, btn)
	}

	if err := cb.buttonProvider.SetButtonGroup(ctx, groupID, buttons); err != nil {
		return nil, "", fmt.Errorf("store buttons group: %w", err)
	}

	return buttons, groupID, nil
}
