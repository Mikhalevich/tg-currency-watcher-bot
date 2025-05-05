package infra

import (
	"flag"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jinzhu/configor"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"

	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/adapter/storage/postgres"
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

	p := postgres.New(dbConn)

	return p, func() {
		dbConn.Close()
	}, nil
}
