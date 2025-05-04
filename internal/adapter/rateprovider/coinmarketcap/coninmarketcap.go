package coinmarketcap

import (
	"net/http"

	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/domain/exchange"
)

var _ exchange.RateProvider = (*CoinMarketCap)(nil)

type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type CoinMarketCap struct {
	apiKey string
	doer   HTTPDoer
}

func New(apiKey string, doer HTTPDoer) *CoinMarketCap {
	return &CoinMarketCap{
		apiKey: apiKey,
		doer:   doer,
	}
}
