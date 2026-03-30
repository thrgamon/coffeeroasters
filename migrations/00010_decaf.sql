-- +goose Up
ALTER TABLE coffees ADD COLUMN is_decaf BOOLEAN NOT NULL DEFAULT false;
UPDATE coffees SET is_decaf = true WHERE lower(name) LIKE '%decaf%' OR lower(name) LIKE '%decaffeinated%';
CREATE INDEX idx_coffees_is_decaf ON coffees(is_decaf) WHERE is_decaf = true;

-- +goose Down
DROP INDEX IF EXISTS idx_coffees_is_decaf;
ALTER TABLE coffees DROP COLUMN is_decaf;
