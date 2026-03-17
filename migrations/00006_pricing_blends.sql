-- +goose Up
-- Per-100g pricing columns
ALTER TABLE coffees ADD COLUMN price_per_100g_min INTEGER;
ALTER TABLE coffees ADD COLUMN price_per_100g_max INTEGER;

-- Blend support
ALTER TABLE coffees ADD COLUMN is_blend BOOLEAN NOT NULL DEFAULT FALSE;

CREATE TABLE blend_components (
    id SERIAL PRIMARY KEY,
    coffee_id INTEGER NOT NULL REFERENCES coffees(id) ON DELETE CASCADE,
    country_code CHAR(2) REFERENCES countries(code),
    region_id INTEGER REFERENCES regions(id),
    variety TEXT,
    percentage INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_blend_components_coffee_id ON blend_components(coffee_id);

-- +goose Down
DROP INDEX IF EXISTS idx_blend_components_coffee_id;
DROP TABLE IF EXISTS blend_components;
ALTER TABLE coffees DROP COLUMN IF EXISTS is_blend;
ALTER TABLE coffees DROP COLUMN IF EXISTS price_per_100g_max;
ALTER TABLE coffees DROP COLUMN IF EXISTS price_per_100g_min;
