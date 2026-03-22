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
SELECT r.id, r.slug, r.name, r.website, r.state,
    count(c.id)::int AS coffee_count
FROM roasters r
LEFT JOIN coffees c ON c.roaster_id = r.id AND c.in_stock = true
WHERE r.active = true AND r.opted_out = false
GROUP BY r.id, r.slug, r.name, r.website, r.state
ORDER BY r.name;

-- name: ListRoastersByState :many
SELECT r.id, r.slug, r.name, r.website, r.state,
    count(c.id)::int AS coffee_count
FROM roasters r
LEFT JOIN coffees c ON c.roaster_id = r.id AND c.in_stock = true
WHERE r.state = $1 AND r.active = true AND r.opted_out = false
GROUP BY r.id, r.slug, r.name, r.website, r.state
ORDER BY r.name;

-- name: CountRoasters :one
SELECT count(*) FROM roasters WHERE active = true AND opted_out = false;
