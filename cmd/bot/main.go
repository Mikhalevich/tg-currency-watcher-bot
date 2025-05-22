package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/adapter/messagesender"
	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/app/currencybot"
	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/config"
	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/domain/button"
	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/domain/rates"
	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/domain/user"
	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/infra"
	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/infra/logger"
	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/infra/tracing"
)

func main() {
	var cfg config.CurrencyBot
	if err := infra.LoadConfig(&cfg); err != nil {
		logger.StdLogger().WithError(err).Error("failed to load config")
		os.Exit(1)
	}

	log, err := infra.SetupLogger(cfg.LogLevel)
	if err != nil {
		logger.StdLogger().WithError(err).Error("failed to setup logger")
		os.Exit(1)
	}

	if err := tracing.SetupTracer(cfg.Tracing.Endpoint, cfg.Tracing.ServiceName, ""); err != nil {
		log.WithError(err).Error("failed to setup tracer")
		os.Exit(1)
	}

	if err := infra.RunSignalInterruptionFunc(func(ctx context.Context) error {
		pDB, cleanup, err := infra.MakePostgres(cfg.Postgres)
		if err != nil {
			return fmt.Errorf("init postgres: %w", err)
		}

		defer cleanup()

		buttonRepository, err := infra.MakeRedisButtonRepository(ctx, cfg.ButtonRedis)
		if err != nil {
			return fmt.Errorf("make redis button repository: %w", err)
		}

		msgSender, err := messagesender.New(cfg.Bot.Token)
		if err != nil {
			return fmt.Errorf("create message sender: %w", err)
		}

		currencyBot, err := currencybot.New(
			cfg.Bot.Token,
			logger.NewLogrus().WithField("bot_name", "currency_bot"),
			user.NewProcessor(pDB, msgSender),
			rates.New(pDB),
			button.NewButtonProvider(buttonRepository),
		)
		if err != nil {
			return fmt.Errorf("create currency bot: %w", err)
		}

		if err := currencyBot.Start(ctx); err != nil {
			return fmt.Errorf("bot start: %w", err)
		}

		return nil
	}); err != nil {
		log.WithError(err).Error("failed run service")
		os.Exit(1)
	}
}
