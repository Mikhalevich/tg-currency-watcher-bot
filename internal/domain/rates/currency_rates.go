package rates

import (
	"context"
	"fmt"
)

func (r *Rates) CurrencyRates(ctx context.Context) ([]Currency, error) {
	currencies, err := r.ratesProvider.GetCurrencyRates(ctx)
	if err != nil {
		return nil, fmt.Errorf("get currencies: %w", err)
	}

	return currencies, nil
}
