package messagesender

import (
	"fmt"

	"github.com/go-telegram/bot"
)

type MessageSender struct {
	botAPI *bot.Bot
}

func New(token string) (*MessageSender, error) {
	botAPI, err := bot.New(
		token,
		bot.WithSkipGetMe(),
	)
	if err != nil {
		return nil, fmt.Errorf("creating bot: %w", err)
	}

	return &MessageSender{
		botAPI: botAPI,
	}, nil
}
