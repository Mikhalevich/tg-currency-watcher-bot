package currencybot

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (cb *CurrencyBot) DefaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	cb.replyTextMessage(ctx, update.Message.Chat.ID, update.Message.ID, "command is not supported")
}
