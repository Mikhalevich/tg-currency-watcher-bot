package config

import "time"

type CurrencyBot struct {
	LogLevel    string      `yaml:"log_level" required:"true"`
	Tracing     Tracing     `yaml:"tracing" required:"true"`
	Bot         Bot         `yaml:"bot" required:"true"`
	Postgres    Postgres    `yaml:"postgres" required:"true"`
	ButtonRedis ButtonRedis `yaml:"button_redis" required:"true"`
}

type Exchange struct {
	LogLevel      string        `yaml:"log_level" required:"true"`
	Tracing       Tracing       `yaml:"tracing" required:"true"`
	CoinMarketCap CoinMarketCap `yaml:"coin_market_cap" required:"true"`
	Postgres      Postgres      `yaml:"postgres" required:"true"`
}

type Notifier struct {
	LogLevel string        `yaml:"log_level" required:"true"`
	Tracing  Tracing       `yaml:"tracing" required:"true"`
	Bot      Bot           `yaml:"bot" required:"true"`
	Postgres Postgres      `yaml:"postgres" required:"true"`
	Interval time.Duration `yaml:"interval" required:"true"`
}

type CoinMarketCap struct {
	APIKey   string        `yaml:"api_key" required:"true"`
	Timeout  time.Duration `yaml:"timeout" required:"true"`
	Interval time.Duration `yaml:"interval" required:"true"`
}

type Bot struct {
	Token string `yaml:"token" required:"true"`
}

type Tracing struct {
	Endpoint    string `yaml:"endpoint" required:"true"`
	ServiceName string `yaml:"service_name" required:"true"`
}

type Postgres struct {
	Connection string `yaml:"connection" required:"true"`
}

type ButtonRedis struct {
	Addr string        `yaml:"addr" required:"true"`
	Pwd  string        `yaml:"pwd" required:"true"`
	DB   int           `yaml:"db" required:"true"`
	TTL  time.Duration `yaml:"ttl" required:"true"`
}
