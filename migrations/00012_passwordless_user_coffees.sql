-- +goose Up

-- Make password_hash optional for passwordless users
ALTER TABLE users ALTER COLUMN password_hash DROP NOT NULL;
ALTER TABLE users ALTER COLUMN password_hash SET DEFAULT '';

-- Admin flag
ALTER TABLE users ADD COLUMN is_admin BOOLEAN NOT NULL DEFAULT false;

-- Magic links for passwordless authentication
CREATE TABLE magic_links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL,
    token VARCHAR(64) UNIQUE NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_magic_links_token ON magic_links (token);
CREATE INDEX idx_magic_links_email ON magic_links (email);

-- User coffee tracking (Letterboxd-style)
-- status: 'wishlist' = want to try, 'logged' = have drunk it
CREATE TABLE user_coffees (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    coffee_id BIGINT NOT NULL REFERENCES coffees(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL CHECK (status IN ('wishlist', 'logged')),
    liked BOOLEAN,
    rating SMALLINT CHECK (rating IS NULL OR (rating >= 1 AND rating <= 5)),
    review TEXT,
    drunk_at DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id, coffee_id)
);

CREATE INDEX idx_user_coffees_user_id ON user_coffees (user_id);
CREATE INDEX idx_user_coffees_coffee_id ON user_coffees (coffee_id);

-- +goose Down

DROP TABLE IF EXISTS user_coffees;
DROP TABLE IF EXISTS magic_links;
ALTER TABLE users DROP COLUMN IF EXISTS is_admin;
ALTER TABLE users ALTER COLUMN password_hash SET NOT NULL;
ALTER TABLE users ALTER COLUMN password_hash DROP DEFAULT;
