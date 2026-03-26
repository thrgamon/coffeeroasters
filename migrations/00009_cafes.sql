-- Cafes: physical cafe locations associated with roasters
CREATE TABLE cafes (
    id          SERIAL PRIMARY KEY,
    roaster_id  INTEGER NOT NULL REFERENCES roasters(id),
    slug        VARCHAR(100) NOT NULL,
    name        VARCHAR(255) NOT NULL,
    address     TEXT,
    suburb      VARCHAR(100),
    state       CHAR(3),
    postcode    VARCHAR(10),
    latitude    DOUBLE PRECISION,
    longitude   DOUBLE PRECISION,
    phone       VARCHAR(50),
    instagram   VARCHAR(255),
    website_url VARCHAR(500),
    image_url   VARCHAR(500),
    active      BOOLEAN DEFAULT true,
    created_at  TIMESTAMPTZ DEFAULT now(),
    updated_at  TIMESTAMPTZ DEFAULT now(),

    UNIQUE (roaster_id, slug)
);

CREATE INDEX idx_cafes_roaster ON cafes(roaster_id);
CREATE INDEX idx_cafes_state ON cafes(state) WHERE state IS NOT NULL;
CREATE INDEX idx_cafes_active ON cafes(active) WHERE active = true;
CREATE INDEX idx_cafes_slug ON cafes(slug);
