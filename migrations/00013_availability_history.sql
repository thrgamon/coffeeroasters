-- +goose Up

-- Track daily availability snapshots for coffees.
-- The scraper inserts a row each run; we keep the latest per day via UNIQUE.
CREATE TABLE IF NOT EXISTS coffee_availability_log (
    id              BIGSERIAL PRIMARY KEY,
    coffee_id       BIGINT NOT NULL REFERENCES coffees(id) ON DELETE CASCADE,
    in_stock        BOOLEAN NOT NULL,
    price_cents     INT,
    recorded_at     DATE NOT NULL DEFAULT CURRENT_DATE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (coffee_id, recorded_at)
);

CREATE INDEX idx_coffee_availability_log_coffee_id ON coffee_availability_log (coffee_id);
CREATE INDEX idx_coffee_availability_log_recorded_at ON coffee_availability_log (recorded_at DESC);

-- +goose Down
DROP TABLE IF EXISTS coffee_availability_log;
