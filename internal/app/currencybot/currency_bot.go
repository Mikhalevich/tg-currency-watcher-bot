package currencybot

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/domain/user"
	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/infra/logger"
)

type UserCurrency interface {
	GetUserCurrencies(ctx context.Context) ([]user.Currency, error)
}

type CurrencyBot struct {
	botAPI *bot.Bot
	logger logger.Logger

	userCurrency UserCurrency
}

func New(
	token string,
	logger logger.Logger,
	userCurrency UserCurrency,
) (*CurrencyBot, error) {
	currencyBot := CurrencyBot{
		logger:       logger,
		userCurrency: userCurrency,
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

func (cb *CurrencyBot) Start(ctx context.Context) error {
	cb.registerHandlers()

	if err := cb.setMyCommands(ctx); err != nil {
		return fmt.Errorf("set my commands: %w", err)
	}

	cb.logger.Info("bot started")

	cb.botAPI.Start(ctx)

	cb.logger.Info("bot stopped")

	return nil
}

func (cb *CurrencyBot) registerHandlers() {
	cb.botAPI.RegisterHandler(
		bot.HandlerTypeMessageText,
		"/my_currencies",
		bot.MatchTypeExact,
		cb.MyCurrencies,
	)
}

func (cb *CurrencyBot) setMyCommands(ctx context.Context) error {
	if _, err := cb.botAPI.SetMyCommands(ctx, &bot.SetMyCommandsParams{
		Commands: []models.BotCommand{
			{
				Command:     "/my_currencies",
				Description: "get all subscribes currencies",
			},
		},
	}); err != nil {
		return fmt.Errorf("set my commands: %w", err)
	}

	return nil
}
