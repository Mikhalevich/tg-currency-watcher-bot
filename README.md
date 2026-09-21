# tg-currency-watcher-bot

[![Build and test](https://github.com/Mikhalevich/tg-currency-watcher-bot/actions/workflows/build_and_test.yml/badge.svg?branch=master)](https://github.com/Mikhalevich/tg-currency-watcher-bot/actions/workflows/build_and_test.yml)
[![golangci-lint](https://github.com/Mikhalevich/tg-currency-watcher-bot/actions/workflows/golangci-lint.yml/badge.svg?branch=master)](https://github.com/Mikhalevich/tg-currency-watcher-bot/actions/workflows/golangci-lint.yml)

A Go microservice application that watches cryptocurrency prices and notifies Telegram users.

## Architecture overview

The project consists of three independent services that share a PostgreSQL database, Redis, and OpenTelemetry tracing via Jaeger:

| Service | Role |
| --- | --- |
| **bot** | Telegram bot that lets users subscribe and unsubscribe to currency pairs and change their notification interval. |
| **exchange** | Background worker that fetches current rates from CoinMarketCap and updates PostgreSQL. |
| **notifier** | Background worker that sends Telegram notifications to users whose notification time is due. |

Infrastructure support:

- **PostgreSQL** for application data and currency rates.
- **Redis** for storing inline Telegram button payloads.
- **Jaeger** for tracing HTTP, PostgreSQL, Redis, and bot handler spans.
- **sql-migrate** for database migrations.
- **sqlboiler** for generated PostgreSQL models.

## Services and behavior

### Bot

The bot uses long polling and registers these exact commands:

- `/subscribed_currencies` — show the user’s subscribed currency pairs.
- `/currency_pairs` — fetch current pairs from the rates provider and display them as subscribe buttons.
- `/change_notification_interval` — display buttons for intervals `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `12`, `16`, `20`, and `24` hours.
- `/unsubscribe` — show subscribed pairs as unsubscribe buttons.

Callback buttons support:

- Subscribe to a currency pair.
- Unsubscribe from a currency pair.
- Change the notification interval.

Domain errors are surfaced as friendly messages:

- `currency already subscribed`
- `No such currency pair`
- `no subscribed currencies`

The bot starts with `bot.WithSkipGetMe()` and `bot.WithDefaultHandler`, so unsupported text messages reply with `command is not supported`.

### Exchange

The exchange service periodically fetches rates from CoinMarketCap using:

- `coin_market_cap.api_key`
- `coin_market_cap.timeout`
- `coin_market_cap.interval`

It stores updated currencies in PostgreSQL. The scheduler runs immediately and then on a fixed `time.Ticker` interval until context cancellation.

### Notifier

The notifier service periodically selects users ready for notifications with a raw SQL query using `next_notification_time < $1` and `FOR UPDATE SKIP LOCKED`. It then sends messages through the Telegram bot API. The run interval and fetch limit are controlled by:

- `notifier.interval`
- `notifier.limit`

## Configuration

All services load configuration with the `configor` library from a YAML file. The binary accepts a `-config` flag; the default is `config/config.yaml`.

Example files are provided in `config/`:

- `config-bot-example.yaml`
- `config-exchange-example.yaml`
- `config-notifier-example.yaml`
- `dbconfig-example.yml`

The `.gitignore` excludes the concrete config files so local examples are not accidentally committed.

### Bot config fields

```yaml
log_level: "info"
tracing:
  endpoint: "jaeger:4318"
  service_name: "currency-watcher-bot"
bot:
  token: "currency-watcher-bot-token"
postgres:
  connection: "host=postgres port=5432 user=bot password=bot dbname=bot sslmode=disable"
button_redis:
  addr: "redis:6379"
  pwd: "redis123"
  db: 1
  ttl: 1h
```

### Exchange config fields

```yaml
log_level: "info"
tracing:
  endpoint: "jaeger:4318"
  service_name: "exchange"
coin_market_cap:
  api_key: "CoinMarketCap api key"
  timeout: 5s
  interval: 5m
postgres:
  connection: "host=postgres port=5432 user=bot password=bot dbname=bot sslmode=disable"
```

### Notifier config fields

```yaml
log_level: "info"
tracing:
  endpoint: "jaeger:4318"
  service_name: "notifer"
postgres:
  connection: "host=postgres port=5432 user=bot password=bot dbname=bot sslmode=disable"
bot:
  token: "currency-watcher-bot-token"
interval: 1m
limit: 100
```

### Postgres connection

The connection string is used with `pgx` through OpenTelemetry-instrumented SQL. The database schema is managed by `sql-migrate`, and generated models are produced by `sqlboiler`.
### Postgres connection

The connection string is used with `pgx` through OpenTelemetry-instrumented SQL. The database schema is managed by `sql-migrate`, and generated models are produced by `sqlboiler`.

## Project layout

```
cmd/
  bot/         # Telegram bot entrypoint
  exchange/    # CoinMarketCap exchange worker entrypoint
  notifier/    # Notification worker entrypoint
config/        # Example YAML configuration files
internal/
  app/currencybot/   # Telegram bot handlers and button logic
  config/            # Config structs and loading
  domain/            # Bot, exchange, user, rates, and button domain models
  adapter/           # Storage, Redis button repository, CoinMarketCap provider, message sender
  infra/             # Logger, tracing, scheduler, Postgres initialization, signal handling
script/
  db/migrations/     # sql-migrate SQL migrations
  docker/            # Dockerfiles and docker-compose files
```

## Database schema

Migrations are in `script/db/migrations/`:

- `1.currency.sql` — creates the `currency` table and seeds popular pairs.
- `2.user.sql` — creates `users`, `users_currency`, and indexes for notification queries.

The generated SQLBoiler models live in `internal/adapter/storage/postgres/internal/models`.

## Building and testing

```bash
# Build all binaries into bin/
make build

# Run all tests
make test

# Install and run golangci-lint
make lint

# Install mockgen
make tools

# Vendor dependencies
make vendor
```

The project uses vendored dependencies and builds with `CGO_ENABLED=0` in Docker.

## Docker and Docker Compose

Start the full local stack:

```bash
make compose-up
```

This starts:

- `bot`
- `exchange`
- `notifier`
- `postgres`
- `redis`
- `jaeger`
- `sql-migrate`

`sql-migrate` runs migrations before the application services start. Jaeger is available on ports `16686`, `14268`, and `4318`.

Stop the stack:

```bash
make compose-down
```

### Docker images

Each service has a multi-stage Dockerfile in `script/docker/`:

- `bot.Dockerfile`
- `exchange.Dockerfile`
- `notifier.Dockerfile`
- `sqlboiler.Dockerfile`
- `sqlmigrate.Dockerfile`

Example entrypoints:

```dockerfile
ENTRYPOINT ["./bot", "-config", "config-bot.yaml"]
ENTRYPOINT ["./exchange", "-config", "config-exchange.yaml"]
ENTRYPOINT ["./notifier", "-config", "config-notifier.yaml"]
```

### Generate SQLBoiler models

```bash
make generate-db-models-up
```

This uses `script/docker/sqlboiler-docker-compose.yml`, which depends on `sql-migrate` and mounts the generated models into `internal/adapter/storage/postgres/internal/models`.

## CI

GitHub Actions runs on push and pull request:

- `.github/workflows/build_and_test.yml` builds and tests the project.
- `.github/workflows/golangci-lint.yml` runs `golangci-lint` against `master`.

Configuration is in `.golangci.yml`.

## Observability

- Logging uses `logrus` in JSON format with OpenTelemetry log hooks.
- Tracing is configured from `tracing.endpoint` and `tracing.service_name`. The bot connects directly to Jaeger via OTLP HTTP.
- The HTTP client used by CoinMarketCap is wrapped with `otelhttp`.
- PostgreSQL is opened with `otelsql`, and Redis is instrumented with `redisotel`.

## License

See the project repository for license details.

