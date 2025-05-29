package infra

import (
	"context"
	"flag"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jinzhu/configor"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"

	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/adapter/buttonrespository"
	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/adapter/storage/postgres"
	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/adapter/storage/postgres/driver"
	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/config"
	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/infra/logger"
)

func LoadConfig(cfg any) error {
	configFile := flag.String("config", "config/config.yaml", "consumer worker config file")
	flag.Parse()

	if err := configor.Load(cfg, *configFile); err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	return nil
}

func SetupLogger(lvl string) (logger.Logger, error) {
	log, err := logger.NewLogrusWithLevel(lvl)
	if err != nil {
		return nil, fmt.Errorf("creating new logger: %w", err)
	}

	logger.SetStdLogger(log)

	return log, nil
}

func MakePostgres(cfg config.Postgres) (*postgres.Postgres, func(), error) {
	if cfg.Connection == "" {
		return nil, func() {}, nil
	}

	dbConn, err := otelsql.Open("pgx", cfg.Connection)
	if err != nil {
		return nil, nil, fmt.Errorf("open database: %w", err)
	}

	if err := dbConn.Ping(); err != nil {
		return nil, nil, fmt.Errorf("ping: %w", err)
	}

	p := postgres.New(dbConn, driver.NewPGX())

	return p, func() {
		dbConn.Close()
	}, nil
}

func MakeRedisButtonRepository(
	ctx context.Context,
	cfg config.ButtonRedis,
) (*buttonrespository.ButtonRepository, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Pwd,
		DB:       cfg.DB,
	})

	if err := redisotel.InstrumentTracing(rdb); err != nil {
		return nil, fmt.Errorf("redis instrument tracing: %w", err)
	}

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}

	return buttonrespository.New(rdb, cfg.TTL), nil
}
