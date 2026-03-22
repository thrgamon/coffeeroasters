-- +goose Up
ALTER TABLE coffees ADD COLUMN description TEXT;
ALTER TABLE coffees ADD COLUMN embedding FLOAT8[];

-- +goose Down
ALTER TABLE coffees DROP COLUMN IF EXISTS embedding;
ALTER TABLE coffees DROP COLUMN IF EXISTS description;
