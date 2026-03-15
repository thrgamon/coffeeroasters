-- Coffeeroasters DB Schema

CREATE TABLE IF NOT EXISTS roasters (
    id          SERIAL PRIMARY KEY,
    name        TEXT NOT NULL,
    slug        TEXT NOT NULL UNIQUE,
    website     TEXT,
    country     TEXT,
    region      TEXT,
    description TEXT,
    logo_url    TEXT,
    source_url  TEXT,
    discovered_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS coffees (
    id          SERIAL PRIMARY KEY,
    roaster_id  INTEGER NOT NULL REFERENCES roasters(id) ON DELETE CASCADE,
    name        TEXT NOT NULL,
    origin      TEXT,
    process     TEXT,
    variety     TEXT,
    roast_level TEXT,
    tasting_notes TEXT[],
    price_cents INTEGER,
    currency    TEXT DEFAULT 'USD',
    weight_grams INTEGER,
    url         TEXT,
    in_stock    BOOLEAN DEFAULT TRUE,
    scraped_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS scrape_runs (
    id          SERIAL PRIMARY KEY,
    roaster_id  INTEGER REFERENCES roasters(id),
    started_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ,
    status      TEXT NOT NULL DEFAULT 'running', -- running, success, failed
    error_msg   TEXT,
    items_found INTEGER DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_coffees_roaster_id ON coffees(roaster_id);
CREATE INDEX IF NOT EXISTS idx_roasters_slug ON roasters(slug);
