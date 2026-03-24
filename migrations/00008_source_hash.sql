-- +goose Up
ALTER TABLE coffees ADD COLUMN source_hash TEXT;

-- +goose Down
ALTER TABLE coffees DROP COLUMN IF EXISTS source_hash;
