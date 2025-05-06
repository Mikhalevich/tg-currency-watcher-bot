package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/app/currencybot"
	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/config"
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

		currencyBot, err := currencybot.New(
			cfg.Bot.Token,
			logger.NewLogrus().WithField("bot_name", "currency_bot"),
			user.New(pDB),
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
