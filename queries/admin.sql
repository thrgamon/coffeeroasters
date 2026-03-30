-- Admin queries for managing roasters, coffees, and cafes.

-- name: AdminListRoasters :many
SELECT id, slug, name, website, state, logo_url, description, active, opted_out, created_at, updated_at
FROM roasters
ORDER BY name;

-- name: AdminGetRoaster :one
SELECT id, slug, name, website, state, logo_url, description, active, opted_out, created_at, updated_at
FROM roasters
WHERE id = $1;

-- name: AdminCreateRoaster :one
INSERT INTO roasters (slug, name, website, state, description, logo_url, active)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id;

-- name: AdminUpdateRoaster :exec
UPDATE roasters SET
    slug = $2,
    name = $3,
    website = $4,
    state = $5,
    description = $6,
    logo_url = $7,
    active = $8,
    opted_out = $9,
    updated_at = now()
WHERE id = $1;

-- name: AdminListCoffees :many
SELECT
    c.id, c.roaster_id, c.name, c.product_url, c.image_url,
    c.country_code, c.process, c.roast_level,
    c.tasting_notes, c.price_cents, c.weight_grams, c.in_stock,
    c.variety, c.species,
    c.price_per_100g_min, c.price_per_100g_max, c.is_blend, c.is_decaf,
    c.description, c.first_seen_at, c.last_seen_at,
    r.name AS roaster_name, r.slug AS roaster_slug,
    co.name AS country_name,
    COUNT(*) OVER() AS total_count
FROM coffees c
JOIN roasters r ON r.id = c.roaster_id
LEFT JOIN countries co ON co.code = c.country_code
ORDER BY c.updated_at DESC
LIMIT sqlc.arg('lim') OFFSET sqlc.arg('off');

-- name: AdminCreateCoffee :one
INSERT INTO coffees (
    roaster_id, name, product_url, image_url,
    country_code, region_id, producer_id,
    process, roast_level, tasting_notes,
    price_cents, weight_grams, price_per_100g_min, price_per_100g_max,
    variety, species, is_blend, is_decaf, in_stock, description, currency
)
VALUES (
    $1, $2, $3, $4,
    $5, $6, $7,
    $8, $9, $10,
    $11, $12, $13, $14,
    $15, $16, $17, $18, $19, $20, 'AUD'
)
RETURNING id;

-- name: AdminUpdateCoffee :exec
UPDATE coffees SET
    roaster_id = $2,
    name = $3,
    product_url = $4,
    image_url = $5,
    country_code = $6,
    region_id = $7,
    producer_id = $8,
    process = $9,
    roast_level = $10,
    tasting_notes = $11,
    price_cents = $12,
    weight_grams = $13,
    price_per_100g_min = $14,
    price_per_100g_max = $15,
    variety = $16,
    species = $17,
    is_blend = $18,
    is_decaf = $19,
    in_stock = $20,
    description = $21,
    updated_at = now()
WHERE id = $1;

-- name: AdminListCafes :many
SELECT c.id, c.roaster_id, c.slug, c.name, c.type, c.address, c.suburb, c.state, c.postcode,
       c.latitude, c.longitude, c.phone, c.instagram, c.website_url, c.image_url, c.active,
       r.name AS roaster_name, r.slug AS roaster_slug
FROM cafes c
JOIN roasters r ON r.id = c.roaster_id
ORDER BY r.name, c.name;

-- name: AdminGetCafe :one
SELECT c.id, c.roaster_id, c.slug, c.name, c.type, c.address, c.suburb, c.state, c.postcode,
       c.latitude, c.longitude, c.phone, c.instagram, c.website_url, c.image_url, c.active,
       r.name AS roaster_name, r.slug AS roaster_slug
FROM cafes c
JOIN roasters r ON r.id = c.roaster_id
WHERE c.id = $1;

-- name: AdminCreateCafe :one
INSERT INTO cafes (roaster_id, slug, name, type, address, suburb, state, postcode,
                   latitude, longitude, phone, instagram, website_url, image_url)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
RETURNING id;

-- name: AdminUpdateCafe :exec
UPDATE cafes SET
    roaster_id = $2,
    slug = $3,
    name = $4,
    type = $5,
    address = $6,
    suburb = $7,
    state = $8,
    postcode = $9,
    latitude = $10,
    longitude = $11,
    phone = $12,
    instagram = $13,
    website_url = $14,
    image_url = $15,
    active = $16,
    updated_at = now()
WHERE id = $1;

-- name: AdminListScrapeRuns :many
SELECT sr.id, sr.roaster_id, sr.started_at, sr.finished_at, sr.status,
       sr.coffees_found, sr.coffees_added, sr.coffees_updated, sr.coffees_removed,
       sr.error_message, sr.duration_ms,
       r.name AS roaster_name, r.slug AS roaster_slug
FROM scrape_runs sr
JOIN roasters r ON r.id = sr.roaster_id
ORDER BY sr.started_at DESC
LIMIT sqlc.arg('lim') OFFSET sqlc.arg('off');
