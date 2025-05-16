package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/adapter/rateprovider/coinmarketcap"
	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/config"
	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/domain/exchange"
	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/infra"
	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/infra/logger"
	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/infra/scheduler"
	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/infra/tracing"
)

func main() {
	var cfg config.Exchange
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
		var (
			httpClient    = tracing.NewClient(cfg.CoinMarketCap.Timeout)
			coinMarketCap = coinmarketcap.New(cfg.CoinMarketCap.APIKey, httpClient)
		)

		pDB, cleanup, err := infra.MakePostgres(cfg.Postgres)
		if err != nil {
			return fmt.Errorf("init postgres: %w", err)
		}

		defer cleanup()

		cmcExchange := exchange.New(pDB, coinMarketCap)

		scheduler.PeriodicTaskExecutor(
			ctx,
			cfg.CoinMarketCap.Interval,
			"coinmarketcap_exchange",
			func(ctx context.Context) error {
				if err := cmcExchange.UpdateCurrencies(ctx); err != nil {
					return fmt.Errorf("update currencies: %w", err)
				}

				return nil
			},
		)

		return nil
	}); err != nil {
		log.WithError(err).Error("failed run service")
		os.Exit(1)
	}
}
