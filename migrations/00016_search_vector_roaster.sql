-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION coffees_search_vector_update() RETURNS trigger AS $$
DECLARE
    roaster_name TEXT;
BEGIN
    SELECT name INTO roaster_name FROM roasters WHERE id = NEW.roaster_id;
    NEW.search_vector :=
        setweight(to_tsvector('english', coalesce(NEW.name, '')), 'A') ||
        setweight(to_tsvector('english', coalesce(roaster_name, '')), 'A') ||
        setweight(to_tsvector('english', coalesce(NEW.origin_raw, '')), 'B') ||
        setweight(to_tsvector('english', coalesce(NEW.tasting_notes_raw, '')), 'C') ||
        setweight(to_tsvector('english', coalesce(NEW.variety_raw, '')), 'D');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

UPDATE coffees SET updated_at = updated_at;

-- +goose Down
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

UPDATE coffees SET updated_at = updated_at;
