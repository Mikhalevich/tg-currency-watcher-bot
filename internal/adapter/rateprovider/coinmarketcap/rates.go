package coinmarketcap

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/adapter/rateprovider/coinmarketcap/internal/contract"
	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/domain/exchange"
)

const (
	baseURL = "https://pro-api.coinmarketcap.com/v2/cryptocurrency/quotes/latest"
	comma   = ","
)

func (c *CoinMarketCap) Rates(
	ctx context.Context,
	convertFrom []exchange.ExternalID,
	convertTo exchange.ExternalID,
) (map[exchange.ExternalID]exchange.Money, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create http request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	//nolint:canonicalheader
	req.Header.Set("X-CMC_PRO_API_KEY", c.apiKey)

	req.URL.RawQuery = makeQuery(convertFrom, convertTo)

	rsp, err := c.doer.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do http request: %w", err)
	}

	defer rsp.Body.Close()

	var quotes contract.QuoteResponse
	if err := json.NewDecoder(rsp.Body).Decode(&quotes); err != nil {
		return nil, fmt.Errorf("json decode: %w", err)
	}

	if quotes.IsError() {
		return nil, fmt.Errorf("response error: %s", quotes.ErrorMessage())
	}

	return quotes.Quotes(convertTo), nil
}

func makeQuery(from []exchange.ExternalID, to exchange.ExternalID) string {
	q := make(url.Values)
	q.Add("convert_id", to.String())
	q.Add("id", makeCommaSeparatedExternalIDs(from))

	return q.Encode()
}

func makeCommaSeparatedExternalIDs(ids []exchange.ExternalID) string {
	switch len(ids) {
	case 0:
		return ""
	case 1:
		return ids[0].String()
	}

	var bufLen int
	bufLen += len(comma)*len(ids) - 1

	for _, id := range ids {
		bufLen += len(id)
	}

	var builder strings.Builder

	builder.Grow(bufLen)

	builder.WriteString(ids[0].String())

	for _, id := range ids[1:] {
		builder.WriteString(comma)
		builder.WriteString(id.String())
	}

	return builder.String()
}
