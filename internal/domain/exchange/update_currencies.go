package exchange

import (
	"context"
	"fmt"
	"time"
)

func (e *Exchange) UpdateCurrencies(ctx context.Context) error {
	currencies, err := e.storage.GetCurrencies(ctx)
	if err != nil {
		return fmt.Errorf("get currencies: %w", err)
	}

	from, to := extractQuotes(currencies)

	rates, err := e.rateProvider.Rates(ctx, from, to)
	if err != nil {
		return fmt.Errorf("get rates: %w", err)
	}

	now := time.Now()

	for i, v := range currencies {
		if price, ok := rates[v.QuoteExternalID]; ok {
			currencies[i].Price = price
			currencies[i].UpdatedAt = now
		}
	}

	if err := e.storage.UpdateCurrencies(ctx, currencies); err != nil {
		return fmt.Errorf("update currencies: %w", err)
	}

	return nil
}

func extractQuotes(currencies []Currency) ([]ExternalID, ExternalID) {
	if len(currencies) == 0 {
		return nil, ""
	}

	quotes := make([]ExternalID, 0, len(currencies))
	for _, c := range currencies {
		quotes = append(quotes, c.QuoteExternalID)
	}

	return quotes, currencies[0].BaseExternalID
}
