package user

import (
	"context"
	"fmt"
	"strings"
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

func (c Currency) FormatPrice() string {
	return strings.TrimSpace(fmt.Sprintf("%9.2f", c.pairPrice()))
}

func (c Currency) pairPrice() Money {
	if c.IsInverted {
		return 1 / c.Price
	}

	return c.Price
}

func (u *User) GetUserCurrencies(ctx context.Context) ([]Currency, error) {
	currencies, err := u.storage.GetUserCurrencies(ctx)
	if err != nil {
		return nil, fmt.Errorf("get user currencies: %w", err)
	}

	return currencies, nil
}
