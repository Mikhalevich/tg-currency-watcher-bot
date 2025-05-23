package currencybot

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (cb *CurrencyBot) DefaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	info, err := textMessageInfo(update)
	if err != nil {
		// skip error
		return
	}

	cb.replyTextMessage(ctx, info.ChatID, info.MessageID, "command is not supported")
}
