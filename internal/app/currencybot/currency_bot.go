package currencybot

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/domain/user"
	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/infra/logger"
	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/infra/tracing"
)

type UserCurrency interface {
	GetUserCurrencies(ctx context.Context) ([]user.Currency, error)
}

type CurrencyBot struct {
	botAPI *bot.Bot
	logger logger.Logger

	userCurrency UserCurrency

	commands []models.BotCommand
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
	cb.registerCommandTextHandler("/my_currencies", "get all subscribed currency pairs", cb.MyCurrencies)
}

func (cb *CurrencyBot) setMyCommands(ctx context.Context) error {
	if cb.commands == nil {
		return nil
	}

	if _, err := cb.botAPI.SetMyCommands(ctx, &bot.SetMyCommandsParams{
		Commands: cb.commands,
	}); err != nil {
		return fmt.Errorf("set my commands: %w", err)
	}

	return nil
}

type botHandler func(ctx context.Context, botAPI *bot.Bot, update *models.Update) error

func (cb *CurrencyBot) registerCommandTextHandler(
	pattern string,
	description string,
	handler botHandler,
) {
	wrapper := func(ctx context.Context, botAPI *bot.Bot, update *models.Update) {
		ctx, span := tracing.StartSpanName(ctx, pattern)
		defer span.End()

		ctx = logger.WithLogger(ctx, cb.logger.WithField("handler_path", pattern))

		if err := handler(ctx, botAPI, update); err != nil {
			cb.replyTextMessage(ctx, update.Message.Chat.ID, update.Message.ID, "internal error")
		}
	}

	cb.botAPI.RegisterHandler(
		bot.HandlerTypeMessageText,
		pattern,
		bot.MatchTypeExact,
		wrapper,
	)

	cb.commands = append(cb.commands, models.BotCommand{
		Command:     pattern,
		Description: description,
	})
}

func (cb *CurrencyBot) replyTextMessage(ctx context.Context, chatID int64, messageID int, text string) {
	if _, err := cb.botAPI.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		ReplyParameters: &models.ReplyParameters{
			MessageID: messageID,
		},
		Text: text,
	}); err != nil {
		cb.logger.WithError(err).Error("send message")
	}
}
