package messagesender

import (
	"context"

	"github.com/go-telegram/bot"

	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/infra/logger"
)

func (m *MessageSender) SendTextMessage(ctx context.Context, chatID int64, text string) {
	if _, err := m.botAPI.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   text,
	}); err != nil {
		logger.FromContext(ctx).WithError(err).Error("send message")
	}
}
