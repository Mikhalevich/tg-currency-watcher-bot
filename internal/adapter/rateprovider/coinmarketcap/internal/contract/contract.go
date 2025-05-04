package contract

import (
	"fmt"

	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/domain/exchange"
)

type Quote struct {
	Price float64 `json:"price"`
}

type Data struct {
	ID     int              `json:"id"`
	Name   string           `json:"name"`
	Symbol string           `json:"symbol"`
	Slug   string           `json:"slug"`
	Quotes map[string]Quote `json:"quote"`
}

type Status struct {
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	CreditCount  int    `json:"credit_count"`
}

type QuoteResponse struct {
	Status Status          `json:"status"`
	Data   map[string]Data `json:"data"`
}

func (q *QuoteResponse) Quotes(toExternalID exchange.ExternalID) map[exchange.ExternalID]exchange.Money {
	if len(q.Data) == 0 {
		return nil
	}

	quotes := make(map[exchange.ExternalID]exchange.Money, len(q.Data))

	for extID, data := range q.Data {
		quotes[exchange.ExternalID(extID)] = quotePrice(data, toExternalID)
	}

	return quotes
}

func quotePrice(data Data, toExternalID exchange.ExternalID) exchange.Money {
	for extID, quote := range data.Quotes {
		if extID == toExternalID.String() {
			if quote.Price > 0 {
				return exchange.Money(1 / quote.Price)
			}

			return 0
		}
	}

	return 0
}

func (q *QuoteResponse) IsError() bool {
	return q.Status.ErrorCode != 0
}

func (q *QuoteResponse) ErrorMessage() string {
	return fmt.Sprintf("code: %d: msg: %s", q.Status.ErrorCode, q.Status.ErrorMessage)
}
