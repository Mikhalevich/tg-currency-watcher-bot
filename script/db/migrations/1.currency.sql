-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE currency(
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    base TEXT NOT NULL,
    base_external_id TEXT NOT NULL,
    quote TEXT NOT NULL,
    quote_external_id TEXT NOT NULL,
    price DECIMAL NOT NULL,
    is_inverted BOOLEAN NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE UNIQUE INDEX base_ext_id_quote_ext_id_idx ON currency(base_external_id, quote_external_id);

INSERT INTO currency(base, base_external_id, quote, quote_external_id, price, is_inverted, updated_at) VALUES
    ('USD', '2781', 'BTC', '1', 0, TRUE, CURRENT_TIMESTAMP),
    ('USD', '2781', 'EUR', '2790', 0, FALSE, CURRENT_TIMESTAMP),
    ('USD', '2781', 'RUB', '2806', 0, FALSE, CURRENT_TIMESTAMP),
    ('USD', '2781', 'BYN', '3533', 0, FALSE, CURRENT_TIMESTAMP),
    ('USD', '2781', 'XAU', '3575', 0, TRUE, CURRENT_TIMESTAMP),
    ('USD', '2781', 'XAG', '3574', 0, TRUE, CURRENT_TIMESTAMP),
    ('USD', '2781', 'XPT', '3577', 0, TRUE, CURRENT_TIMESTAMP),
    ('USD', '2781', 'XPD', '3576', 0, TRUE, CURRENT_TIMESTAMP);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP INDEX base_ext_id_quote_ext_id_idx;
DROP TABLE currency;
