package currencybot

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"

	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/infra/logger"
)

type CurrencyBot struct {
	botAPI *bot.Bot
	logger logger.Logger
}

func New(token string, logger logger.Logger) (*CurrencyBot, error) {
	currencyBot := CurrencyBot{
		logger: logger,
	}

	botAPI, err := bot.New(
		token,
		bot.WithSkipGetMe(),
		bot.WithDefaultHandler(currencyBot.DefaultHandler),
	)
	if err != nil {
		return nil, fmt.Errorf("creating bot: %w", err)
	}

	currencyBot.botAPI = botAPI

	return &currencyBot, nil
}

func (cb *CurrencyBot) Start(ctx context.Context) {
	cb.logger.Info("bot started")

	cb.botAPI.Start(ctx)

	cb.logger.Info("bot stopped")
}
