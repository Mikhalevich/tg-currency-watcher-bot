package user

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

type ExternalID string

func (e ExternalID) String() string {
	return string(e)
}

type Currency struct {
	Base      Symbol
	Quote     Symbol
	Price     Money
	UpdatedAt time.Time
}

func (u *User) GetUserCurrencies(ctx context.Context) ([]Currency, error) {
	currencies, err := u.storage.GetUserCurrencies(ctx)
	if err != nil {
		return nil, fmt.Errorf("get user currencies: %w", err)
	}

	return currencies, nil
}
