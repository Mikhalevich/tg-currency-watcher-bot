package config

import "time"

type CurrencyBot struct {
	LogLevel string  `yaml:"log_level" required:"true"`
	Tracing  Tracing `yaml:"tracing" required:"true"`
	Bot      Bot     `yaml:"bot" required:"true"`
}

type Exchange struct {
	LogLevel      string        `yaml:"log_level" required:"true"`
	Tracing       Tracing       `yaml:"tracing" required:"true"`
	CoinMarketCap CoinMarketCap `yaml:"coin_market_cap" required:"true"`
}

type CoinMarketCap struct {
	APIKey  string        `yaml:"api_key" required:"true"`
	Timeout time.Duration `yaml:"timeout" required:"true"`
}

type Bot struct {
	Token string `yaml:"token" required:"true"`
}

type Tracing struct {
	Endpoint    string `yaml:"endpoint" required:"true"`
	ServiceName string `yaml:"service_name" required:"true"`
}
