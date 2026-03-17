-- +goose Up

-- roasters: the list of Australian indie coffee roasters we track.
CREATE TABLE roasters (
    id          SERIAL PRIMARY KEY,
    slug        VARCHAR(100) UNIQUE NOT NULL,
    name        VARCHAR(255) NOT NULL,
    website     VARCHAR(500) NOT NULL,
    state       CHAR(3),                          -- AU state/territory code
    description TEXT,
    opted_out   BOOLEAN NOT NULL DEFAULT false,   -- roaster has requested removal
    active      BOOLEAN NOT NULL DEFAULT true,    -- included in scheduled scrapes
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_roasters_slug    ON roasters (slug);
CREATE INDEX idx_roasters_state   ON roasters (state);
CREATE INDEX idx_roasters_active  ON roasters (active) WHERE active = true;

-- coffees: individual coffee products scraped from roaster websites.
-- Raw fields store exactly what was scraped; normalised fields are derived.
CREATE TABLE coffees (
    id          BIGSERIAL PRIMARY KEY,
    roaster_id  INTEGER NOT NULL REFERENCES roasters(id) ON DELETE CASCADE,

    -- Identity
    name        VARCHAR(500) NOT NULL,
    product_url VARCHAR(1000),
    image_url   VARCHAR(1000),

    -- Raw scraped values (preserved exactly as found on the site)
    origin_raw       TEXT,
    region_raw       TEXT,
    variety_raw      TEXT,
    process_raw      TEXT,
    roast_raw        TEXT,
    tasting_notes_raw TEXT,
    price_raw        TEXT,
    weight_raw       TEXT,
    currency         CHAR(3) NOT NULL DEFAULT 'AUD',
    in_stock         BOOLEAN NOT NULL DEFAULT true,

    -- Normalised values (derived by internal/normalise)
    -- Origin normalised to country code (ISO 3166-1 alpha-2)
    origin_country   CHAR(2),        -- e.g. 'ET' (Ethiopia)
    origin_region    VARCHAR(255),   -- e.g. 'Yirgacheffe'
    process          VARCHAR(50),    -- canonical: washed|natural|honey|anaerobic|wet-hulled|experimental
    roast_level      VARCHAR(20),    -- canonical: light|medium-light|medium|medium-dark|dark
    tasting_notes    TEXT[],         -- normalised array of individual notes
    price_cents      INTEGER,        -- price in AUD cents
    weight_grams     INTEGER,        -- weight in grams

    -- Full-text search vector (updated by trigger)
    search_vector    tsvector,

    -- Timestamps
    first_seen_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_seen_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_changed_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now(),

    -- Deduplication: one row per roaster+name combination (latest scrape wins)
    UNIQUE (roaster_id, name)
);

CREATE INDEX idx_coffees_roaster_id    ON coffees (roaster_id);
CREATE INDEX idx_coffees_origin        ON coffees (origin_country);
CREATE INDEX idx_coffees_process       ON coffees (process);
CREATE INDEX idx_coffees_roast_level   ON coffees (roast_level);
CREATE INDEX idx_coffees_in_stock      ON coffees (in_stock) WHERE in_stock = true;
CREATE INDEX idx_coffees_search        ON coffees USING gin(search_vector);

-- scrape_runs: audit log of every scrape run (success or failure).
CREATE TABLE scrape_runs (
    id           BIGSERIAL PRIMARY KEY,
    roaster_id   INTEGER NOT NULL REFERENCES roasters(id) ON DELETE CASCADE,
    started_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    finished_at  TIMESTAMPTZ,
    status       VARCHAR(20) NOT NULL DEFAULT 'running',  -- running|success|failed
    coffees_found    INTEGER,
    coffees_added    INTEGER,
    coffees_updated  INTEGER,
    coffees_removed  INTEGER,
    pages_visited    INTEGER,
    error_message    TEXT,
    duration_ms      INTEGER
);

CREATE INDEX idx_scrape_runs_roaster   ON scrape_runs (roaster_id);
CREATE INDEX idx_scrape_runs_status    ON scrape_runs (status);
CREATE INDEX idx_scrape_runs_started   ON scrape_runs (started_at DESC);

-- Trigger: update search_vector on coffees insert/update.
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION coffees_search_vector_update() RETURNS trigger AS $$
BEGIN
    NEW.search_vector :=
        setweight(to_tsvector('english', coalesce(NEW.name, '')), 'A') ||
        setweight(to_tsvector('english', coalesce(NEW.origin_raw, '')), 'B') ||
        setweight(to_tsvector('english', coalesce(NEW.tasting_notes_raw, '')), 'C') ||
        setweight(to_tsvector('english', coalesce(NEW.variety_raw, '')), 'D');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER coffees_search_vector_trigger
    BEFORE INSERT OR UPDATE ON coffees
    FOR EACH ROW EXECUTE FUNCTION coffees_search_vector_update();

-- +goose Down

DROP TRIGGER IF EXISTS coffees_search_vector_trigger ON coffees;
DROP FUNCTION IF EXISTS coffees_search_vector_update();
DROP TABLE IF EXISTS scrape_runs;
DROP TABLE IF EXISTS coffees;
DROP TABLE IF EXISTS roasters;
