-- +goose Up

-- Store brew recipe/notes scraped from roaster sites
ALTER TABLE coffees ADD COLUMN brew_recipe_raw TEXT;

-- User-created brew recipes for specific coffees
CREATE TABLE brew_recipes (
    id          SERIAL PRIMARY KEY,
    coffee_id   BIGINT NOT NULL REFERENCES coffees(id) ON DELETE CASCADE,
    user_id     INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title       VARCHAR(200) NOT NULL,
    brew_method VARCHAR(50) NOT NULL,  -- e.g. espresso, pourover, aeropress, french_press, cold_brew, filter, moka_pot
    dose_grams  NUMERIC(5,1),          -- coffee dose in grams
    water_ml    INTEGER,               -- water volume in ml
    water_temp_c INTEGER,              -- water temperature in celsius
    grind_size  VARCHAR(30),           -- e.g. fine, medium-fine, medium, medium-coarse, coarse
    brew_time_seconds INTEGER,         -- total brew time
    notes       TEXT,                  -- free-form notes, tips, steps
    is_public   BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_brew_recipes_coffee_id ON brew_recipes(coffee_id);
CREATE INDEX idx_brew_recipes_user_id ON brew_recipes(user_id);

-- +goose Down
DROP TABLE IF EXISTS brew_recipes;
ALTER TABLE coffees DROP COLUMN IF EXISTS brew_recipe_raw;
