-- +goose Up
ALTER TABLE coffees ADD COLUMN variety TEXT;
ALTER TABLE coffees ADD COLUMN species TEXT;

-- +goose Down
ALTER TABLE coffees DROP COLUMN variety;
ALTER TABLE coffees DROP COLUMN species;
