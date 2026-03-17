-- +goose Up
ALTER TABLE regions ADD COLUMN latitude DOUBLE PRECISION;
ALTER TABLE regions ADD COLUMN longitude DOUBLE PRECISION;

-- +goose Down
ALTER TABLE regions DROP COLUMN latitude;
ALTER TABLE regions DROP COLUMN longitude;
