-- name: UpsertRoaster :one
INSERT INTO roasters (slug, name, website, state)
VALUES ($1, $2, $3, $4)
ON CONFLICT (slug) DO UPDATE SET
    name = EXCLUDED.name,
    website = EXCLUDED.website,
    state = EXCLUDED.state,
    updated_at = now()
RETURNING id;

-- name: GetRoasterBySlug :one
SELECT id, slug, name, website, state, description, active, created_at, updated_at
FROM roasters
WHERE slug = $1 AND opted_out = false;

-- name: ListRoasters :many
SELECT id, slug, name, website, state
FROM roasters
WHERE active = true AND opted_out = false
ORDER BY name;

-- name: ListRoastersByState :many
SELECT id, slug, name, website, state
FROM roasters
WHERE state = $1 AND active = true AND opted_out = false
ORDER BY name;

-- name: CountRoasters :one
SELECT count(*) FROM roasters WHERE active = true AND opted_out = false;
