-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE currency(
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    base TEXT NOT NULL,
    base_external_id TEXT NOT NULL,
    quote TEXT NOT NULL,
    quote_external_id TEXT NOT NULL,
    price DECIMAL NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE UNIQUE INDEX base_ext_id_quote_ext_id_idx ON currency(base_external_id, quote_external_id);

INSERT INTO currency(base, base_external_id, quote, quote_external_id, price, updated_at) VALUES
    ('USD', '2781', 'BTC', '1', 0, CURRENT_TIMESTAMP);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP INDEX base_ext_id_quote_ext_id_idx;
DROP TABLE currency;
