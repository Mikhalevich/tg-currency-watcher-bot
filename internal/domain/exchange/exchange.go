package exchange

import (
	"context"
	"time"
)

type Money float64
type Symbol string

func (s Symbol) String() string {
	return string(s)
}

type ExternalID string

func (e ExternalID) String() string {
	return string(e)
}

type Currency struct {
	ID              int
	Base            Symbol
	BaseExternalID  ExternalID
	Quote           Symbol
	QuoteExternalID ExternalID
	Price           Money
	UpdatedAt       time.Time
}

type Storage interface {
	GetCurrencies(ctx context.Context) ([]Currency, error)
	UpdateCurrencies(ctx context.Context, currencies []Currency) error
}

type RateProvider interface {
	Rates(ctx context.Context, from []ExternalID, to ExternalID) (map[ExternalID]Money, error)
}

type Exchange struct {
	storage      Storage
	rateProvider RateProvider
}

func New(storage Storage, rateProvider RateProvider) *Exchange {
	return &Exchange{
		storage:      storage,
		rateProvider: rateProvider,
	}
}
