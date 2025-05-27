package currencybot

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/domain/button"
	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/domain/rates"
	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/domain/user"
	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/infra/logger"
	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/infra/tracing"
)

type UserCurrency interface {
	GetCurrenciesByChatID(ctx context.Context, chatID int64) ([]user.Currency, error)
	SubscribeCurrency(ctx context.Context, chatID int64, currencyID int) error
	UnsubscribeCurrency(ctx context.Context, chatID int64, currencyID int) error
	ChangeNotificationInterval(ctx context.Context, chatID int64, interval int) error
}

type RatesProvider interface {
	CurrencyRates(ctx context.Context) ([]rates.Currency, error)
}

type ButtonProvider interface {
	GetButton(ctx context.Context, btnID string) (*button.Button, error)
	SetButtonGroup(ctx context.Context, groupID string, buttons []button.Button) error
}

type CurrencyBot struct {
	botAPI *bot.Bot
	logger logger.Logger

	userCurrency   UserCurrency
	ratesProvider  RatesProvider
	buttonProvider ButtonProvider

	commands []models.BotCommand
}

func New(
	token string,
	logger logger.Logger,
	userCurrency UserCurrency,
	ratesProvider RatesProvider,
	buttonProvider ButtonProvider,
) (*CurrencyBot, error) {
	currencyBot := CurrencyBot{
		logger:         logger,
		userCurrency:   userCurrency,
		ratesProvider:  ratesProvider,
		buttonProvider: buttonProvider,
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
	cb.registerCommandTextHandler("/subscribed_currencies", "get all subscribed currency pairs", cb.MyCurrencies)
	cb.registerCommandTextHandler("/currency_pairs", "view all currency pairs to subscribe", cb.CurrencyPairs)
	cb.registerCommandTextHandler("/change_notification_interval", "change notification interval", cb.NotificationInterval)
	cb.registerCommandTextHandler("/unsubscribe", "view all currencies for unsubscribe", cb.UnsubscribeCurrencies)
	cb.addDefaultCallbackQueryHander(cb.DefaultCallbackQueryHandler)
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

type MessageInfo struct {
	ChatID    int64
	MessageID int
	Data      string
}

type botHandler func(ctx context.Context, botAPI *bot.Bot, info MessageInfo) error
type msgInfoFunc func(update *models.Update) (MessageInfo, error)

func (cb *CurrencyBot) wrapHandler(
	pattern string,
	handler botHandler,
	msgInfoFn msgInfoFunc,
) bot.HandlerFunc {
	return func(ctx context.Context, botAPI *bot.Bot, update *models.Update) {
		ctx, span := tracing.StartSpanName(ctx, pattern)
		defer span.End()

		ctx = logger.WithLogger(ctx, cb.logger.WithField("handler_path", pattern))

		info, err := msgInfoFn(update)
		if err != nil {
			logger.FromContext(ctx).WithError(err).Error("retrieve message info error")

			return
		}

		if err := handler(ctx, botAPI, info); err != nil {
			logger.FromContext(ctx).WithError(err).Error("handler error")

			cb.replyTextMessage(ctx, info.ChatID, info.MessageID, "internal error")
		}
	}
}

func textMessageInfo(update *models.Update) (MessageInfo, error) {
	if update.Message != nil {
		return MessageInfo{
			ChatID:    update.Message.Chat.ID,
			MessageID: update.Message.ID,
		}, nil
	}

	return MessageInfo{}, errors.New("invalid text message")
}

func callbackQueryMessageInfo(update *models.Update) (MessageInfo, error) {
	if update.CallbackQuery != nil {
		return MessageInfo{
			ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
			MessageID: update.CallbackQuery.Message.Message.ID,
			Data:      update.CallbackQuery.Data,
		}, nil
	}

	return MessageInfo{}, errors.New("invalid callback query")
}

func (cb *CurrencyBot) registerCommandTextHandler(
	pattern string,
	description string,
	handler botHandler,
) {
	cb.botAPI.RegisterHandler(
		bot.HandlerTypeMessageText,
		pattern,
		bot.MatchTypeExact,
		cb.wrapHandler(pattern, handler, textMessageInfo),
	)

	cb.commands = append(cb.commands, models.BotCommand{
		Command:     pattern,
		Description: description,
	})
}

func (cb *CurrencyBot) addDefaultCallbackQueryHander(handler botHandler) {
	cb.botAPI.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		"",
		bot.MatchTypePrefix,
		cb.wrapHandler("default_callback_query", handler, callbackQueryMessageInfo),
	)
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

func (cb *CurrencyBot) sendTextMessage(ctx context.Context, chatID int64, text string) {
	if _, err := cb.botAPI.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   text,
	}); err != nil {
		cb.logger.WithError(err).Error("send message")
	}
}
