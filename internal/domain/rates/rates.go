package rates

import (
	"context"
	"time"
)

type Money float64
type Symbol string

func (s Symbol) String() string {
	return string(s)
}

type Currency struct {
	ID        int
	Base      Symbol
	Quote     Symbol
	Price     Money
	UpdatedAt time.Time
}

type RatesProvider interface {
	GetCurrencyRates(ctx context.Context) ([]Currency, error)
}

type Rates struct {
	ratesProvider RatesProvider
}

func New(ratesProvider RatesProvider) *Rates {
	return &Rates{
		ratesProvider: ratesProvider,
	}
}
