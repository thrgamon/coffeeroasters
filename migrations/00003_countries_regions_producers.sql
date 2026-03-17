-- +goose Up

-- countries: seeded reference data for specialty coffee origins.
CREATE TABLE countries (
    code    CHAR(2) PRIMARY KEY,
    name    VARCHAR(100) NOT NULL,
    alpha3  CHAR(3)
);

INSERT INTO countries (code, name, alpha3) VALUES
    ('ET', 'Ethiopia', 'ETH'),
    ('CO', 'Colombia', 'COL'),
    ('BR', 'Brazil', 'BRA'),
    ('KE', 'Kenya', 'KEN'),
    ('GT', 'Guatemala', 'GTM'),
    ('CR', 'Costa Rica', 'CRI'),
    ('PA', 'Panama', 'PAN'),
    ('HN', 'Honduras', 'HND'),
    ('SV', 'El Salvador', 'SLV'),
    ('NI', 'Nicaragua', 'NIC'),
    ('MX', 'Mexico', 'MEX'),
    ('PE', 'Peru', 'PER'),
    ('BO', 'Bolivia', 'BOL'),
    ('EC', 'Ecuador', 'ECU'),
    ('RW', 'Rwanda', 'RWA'),
    ('BI', 'Burundi', 'BDI'),
    ('TZ', 'Tanzania', 'TZA'),
    ('UG', 'Uganda', 'UGA'),
    ('CD', 'DR Congo', 'COD'),
    ('MW', 'Malawi', 'MWI'),
    ('ZM', 'Zambia', 'ZMB'),
    ('ID', 'Indonesia', 'IDN'),
    ('PG', 'Papua New Guinea', 'PNG'),
    ('IN', 'India', 'IND'),
    ('MM', 'Myanmar', 'MMR'),
    ('VN', 'Vietnam', 'VNM'),
    ('CN', 'China', 'CHN'),
    ('TH', 'Thailand', 'THA'),
    ('LA', 'Laos', 'LAO'),
    ('PH', 'Philippines', 'PHL'),
    ('TW', 'Taiwan', 'TWN'),
    ('YE', 'Yemen', 'YEM'),
    ('JM', 'Jamaica', 'JAM'),
    ('HT', 'Haiti', 'HTI'),
    ('DO', 'Dominican Republic', 'DOM'),
    ('CU', 'Cuba', 'CUB');

-- regions: discovered from scraping, get-or-create.
CREATE TABLE regions (
    id           SERIAL PRIMARY KEY,
    country_code CHAR(2) NOT NULL REFERENCES countries(code),
    name         VARCHAR(255) NOT NULL,
    UNIQUE (country_code, name)
);

-- producers: farms/estates/cooperatives, get-or-create.
CREATE TABLE producers (
    id           SERIAL PRIMARY KEY,
    name         VARCHAR(500) NOT NULL,
    country_code CHAR(2) REFERENCES countries(code),
    region_id    INTEGER REFERENCES regions(id),
    UNIQUE (name, country_code)
);

-- Add new FK columns to coffees
ALTER TABLE coffees ADD COLUMN country_code CHAR(2) REFERENCES countries(code);
ALTER TABLE coffees ADD COLUMN region_id    INTEGER REFERENCES regions(id);
ALTER TABLE coffees ADD COLUMN producer_id  INTEGER REFERENCES producers(id);
ALTER TABLE coffees ADD COLUMN producer_raw TEXT;

-- Drop always-NULL columns
ALTER TABLE coffees DROP COLUMN origin_country;
ALTER TABLE coffees DROP COLUMN origin_region;

-- Indexes on the new FK columns
CREATE INDEX idx_coffees_country_code ON coffees (country_code);
CREATE INDEX idx_coffees_region_id    ON coffees (region_id);
CREATE INDEX idx_coffees_producer_id  ON coffees (producer_id);

-- +goose Down

-- Restore dropped columns
ALTER TABLE coffees ADD COLUMN origin_country CHAR(2);
ALTER TABLE coffees ADD COLUMN origin_region  VARCHAR(255);

-- Drop new columns
ALTER TABLE coffees DROP COLUMN IF EXISTS producer_raw;
ALTER TABLE coffees DROP COLUMN IF EXISTS producer_id;
ALTER TABLE coffees DROP COLUMN IF EXISTS region_id;
ALTER TABLE coffees DROP COLUMN IF EXISTS country_code;

DROP INDEX IF EXISTS idx_coffees_producer_id;
DROP INDEX IF EXISTS idx_coffees_region_id;
DROP INDEX IF EXISTS idx_coffees_country_code;

DROP TABLE IF EXISTS producers;
DROP TABLE IF EXISTS regions;
DROP TABLE IF EXISTS countries;
