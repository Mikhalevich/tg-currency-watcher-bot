package rates

import (
	"context"
	"fmt"
	"time"
)

type Money float64
type Symbol string

func (s Symbol) String() string {
	return string(s)
}

type Currency struct {
	ID         int
	Base       Symbol
	Quote      Symbol
	Price      Money
	IsInverted bool
	UpdatedAt  time.Time
}

func (c Currency) FormatPair() string {
	if c.IsInverted {
		return fmt.Sprintf("%s/%s", c.Quote, c.Base)
	}

	return fmt.Sprintf("%s/%s", c.Base, c.Quote)
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
