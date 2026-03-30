-- +goose Up

-- Crowdsourced tasting notes: users vote on existing notes or suggest new ones
CREATE TABLE user_tasting_notes (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    coffee_id BIGINT NOT NULL REFERENCES coffees(id) ON DELETE CASCADE,
    tasting_note VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id, coffee_id, tasting_note)
);

CREATE INDEX idx_user_tasting_notes_coffee ON user_tasting_notes (coffee_id);
CREATE INDEX idx_user_tasting_notes_user ON user_tasting_notes (user_id);

-- +goose Down

DROP TABLE IF EXISTS user_tasting_notes;
