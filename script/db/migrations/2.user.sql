-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE users(
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    chat_id INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    notification_interval_hours INTEGER NOT NULL,
    last_notification_time TIMESTAMPTZ NOT NULL,
    next_notification_time TIMESTAMPTZ GENERATED ALWAYS AS (date_add(last_notification_time, notification_interval_hours * '1 hour'::interval, 'UTC')) STORED
);

CREATE UNIQUE INDEX users_chat_id_idx ON users(chat_id);
CREATE INDEX users_next_notification_time_idx ON users(next_notification_time);

CREATE TABLE users_currency(
    user_id INTEGER NOT NULL,
    currency_id INTEGER NOT NULL,

    CONSTRAINT user_currency_pk PRIMARY KEY(user_id, currency_id),
    CONSTRAINT user_currency_user_id_fk FOREIGN KEY(user_id) REFERENCES users(id),
    CONSTRAINT user_currency_currency_id_fk FOREIGN KEY(currency_id) REFERENCES currency(id)
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE users_currency;
DROP INDEX users_chat_id_idx;
DROP TABLE users;
