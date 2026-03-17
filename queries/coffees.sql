-- name: UpsertCoffee :one
-- Returns true if the row was newly inserted, false if updated.
INSERT INTO coffees (
    roaster_id, name, product_url, image_url,
    origin_raw, region_raw, variety_raw, process_raw, roast_raw,
    tasting_notes_raw, price_raw, weight_raw, currency, in_stock,
    process, roast_level, tasting_notes, price_cents, weight_grams,
    country_code, region_id, producer_id, producer_raw,
    variety, species,
    last_seen_at
)
VALUES (
    $1, $2, $3, $4,
    $5, $6, $7, $8, $9,
    $10, $11, $12, $13, $14,
    $15, $16, $17, $18, $19,
    $20, $21, $22, $23,
    $24, $25,
    now()
)
ON CONFLICT (roaster_id, name) DO UPDATE SET
    product_url = EXCLUDED.product_url,
    image_url = EXCLUDED.image_url,
    origin_raw = EXCLUDED.origin_raw,
    region_raw = EXCLUDED.region_raw,
    variety_raw = EXCLUDED.variety_raw,
    process_raw = EXCLUDED.process_raw,
    roast_raw = EXCLUDED.roast_raw,
    tasting_notes_raw = EXCLUDED.tasting_notes_raw,
    price_raw = EXCLUDED.price_raw,
    weight_raw = EXCLUDED.weight_raw,
    currency = EXCLUDED.currency,
    in_stock = EXCLUDED.in_stock,
    process = EXCLUDED.process,
    roast_level = EXCLUDED.roast_level,
    tasting_notes = EXCLUDED.tasting_notes,
    price_cents = EXCLUDED.price_cents,
    weight_grams = EXCLUDED.weight_grams,
    country_code = EXCLUDED.country_code,
    region_id = EXCLUDED.region_id,
    producer_id = EXCLUDED.producer_id,
    producer_raw = EXCLUDED.producer_raw,
    variety = EXCLUDED.variety,
    species = EXCLUDED.species,
    last_seen_at = now(),
    last_changed_at = CASE
        WHEN coffees.price_cents IS DISTINCT FROM EXCLUDED.price_cents
            OR coffees.in_stock IS DISTINCT FROM EXCLUDED.in_stock
        THEN now()
        ELSE coffees.last_changed_at
    END,
    updated_at = now()
RETURNING (xmax = 0) AS is_new;

-- name: ListCoffees :many
SELECT
    c.id, c.roaster_id, c.name, c.product_url, c.image_url,
    c.country_code, c.origin_raw, c.region_raw, c.process, c.roast_level,
    c.tasting_notes, c.price_cents, c.weight_grams, c.in_stock,
    c.variety, c.species,
    r.name AS roaster_name, r.slug AS roaster_slug,
    co.name AS country_name,
    reg.name AS region_name, reg.id AS coffee_region_id,
    p.name AS producer_name, p.id AS coffee_producer_id
FROM coffees c
JOIN roasters r ON r.id = c.roaster_id
LEFT JOIN countries co ON co.code = c.country_code
LEFT JOIN regions reg ON reg.id = c.region_id
LEFT JOIN producers p ON p.id = c.producer_id
WHERE r.opted_out = false
    AND c.in_stock = true
ORDER BY c.name
LIMIT $1 OFFSET $2;

-- name: CountCoffees :one
SELECT count(*)
FROM coffees c
JOIN roasters r ON r.id = c.roaster_id
WHERE r.opted_out = false
    AND c.in_stock = true;

-- name: ListCoffeesByRoaster :many
SELECT
    c.id, c.roaster_id, c.name, c.product_url, c.image_url,
    c.country_code, c.origin_raw, c.region_raw, c.process, c.roast_level,
    c.tasting_notes, c.price_cents, c.weight_grams, c.in_stock,
    c.variety, c.species,
    r.name AS roaster_name, r.slug AS roaster_slug,
    co.name AS country_name,
    reg.name AS region_name, reg.id AS coffee_region_id,
    p.name AS producer_name, p.id AS coffee_producer_id
FROM coffees c
JOIN roasters r ON r.id = c.roaster_id
LEFT JOIN countries co ON co.code = c.country_code
LEFT JOIN regions reg ON reg.id = c.region_id
LEFT JOIN producers p ON p.id = c.producer_id
WHERE r.slug = $1 AND r.opted_out = false
ORDER BY c.name;

-- name: GetCoffeeByID :one
SELECT
    c.id, c.roaster_id, c.name, c.product_url, c.image_url,
    c.origin_raw, c.region_raw, c.variety_raw, c.process_raw, c.roast_raw,
    c.tasting_notes_raw, c.country_code, c.process, c.roast_level,
    c.tasting_notes, c.price_cents, c.weight_grams, c.in_stock,
    c.producer_raw, c.region_id, c.producer_id,
    c.variety, c.species,
    c.first_seen_at, c.last_seen_at,
    r.name AS roaster_name, r.slug AS roaster_slug,
    co.name AS country_name,
    reg.name AS region_name,
    p.name AS producer_name
FROM coffees c
JOIN roasters r ON r.id = c.roaster_id
LEFT JOIN countries co ON co.code = c.country_code
LEFT JOIN regions reg ON reg.id = c.region_id
LEFT JOIN producers p ON p.id = c.producer_id
WHERE c.id = $1;

-- name: SearchCoffees :many
SELECT
    c.id, c.roaster_id, c.name, c.product_url, c.image_url,
    c.country_code, c.origin_raw, c.region_raw, c.process, c.roast_level,
    c.tasting_notes, c.price_cents, c.weight_grams, c.in_stock,
    c.variety, c.species,
    r.name AS roaster_name, r.slug AS roaster_slug,
    co.name AS country_name,
    reg.name AS region_name, reg.id AS coffee_region_id,
    p.name AS producer_name, p.id AS coffee_producer_id
FROM coffees c
JOIN roasters r ON r.id = c.roaster_id
LEFT JOIN countries co ON co.code = c.country_code
LEFT JOIN regions reg ON reg.id = c.region_id
LEFT JOIN producers p ON p.id = c.producer_id
WHERE r.opted_out = false
    AND c.in_stock = true
    AND c.search_vector @@ plainto_tsquery('english', $1)
ORDER BY ts_rank(c.search_vector, plainto_tsquery('english', $1)) DESC
LIMIT $2 OFFSET $3;

-- name: FilterCoffees :many
SELECT
    c.id, c.roaster_id, c.name, c.product_url, c.image_url,
    c.country_code, c.origin_raw, c.region_raw, c.process, c.roast_level,
    c.tasting_notes, c.price_cents, c.weight_grams, c.in_stock,
    c.variety, c.species,
    r.name AS roaster_name, r.slug AS roaster_slug,
    co.name AS country_name,
    reg.name AS region_name, reg.id AS coffee_region_id,
    p.name AS producer_name, p.id AS coffee_producer_id
FROM coffees c
JOIN roasters r ON r.id = c.roaster_id
LEFT JOIN countries co ON co.code = c.country_code
LEFT JOIN regions reg ON reg.id = c.region_id
LEFT JOIN producers p ON p.id = c.producer_id
WHERE r.opted_out = false
    AND ($1::text IS NULL OR c.country_code = $1)
    AND ($2::text IS NULL OR c.process = $2)
    AND ($3::text IS NULL OR c.roast_level = $3)
    AND ($4::boolean IS NULL OR c.in_stock = $4)
    AND ($5::text IS NULL OR c.variety = $5)
ORDER BY c.name
LIMIT $6 OFFSET $7;

-- name: ListDistinctOrigins :many
SELECT DISTINCT origin_raw
FROM coffees
WHERE origin_raw IS NOT NULL AND origin_raw != ''
ORDER BY origin_raw;

-- name: ListCoffeesByCountry :many
SELECT
    c.id, c.roaster_id, c.name, c.product_url, c.image_url,
    c.country_code, c.process, c.roast_level,
    c.tasting_notes, c.price_cents, c.weight_grams, c.in_stock,
    c.variety, c.species,
    r.name AS roaster_name, r.slug AS roaster_slug,
    co.name AS country_name,
    reg.name AS region_name, reg.id AS coffee_region_id,
    p.name AS producer_name, p.id AS coffee_producer_id
FROM coffees c
JOIN roasters r ON r.id = c.roaster_id
LEFT JOIN countries co ON co.code = c.country_code
LEFT JOIN regions reg ON reg.id = c.region_id
LEFT JOIN producers p ON p.id = c.producer_id
WHERE c.country_code = $1 AND r.opted_out = false
ORDER BY c.name;

-- name: ListCoffeesByRegion :many
SELECT
    c.id, c.roaster_id, c.name, c.product_url, c.image_url,
    c.country_code, c.process, c.roast_level,
    c.tasting_notes, c.price_cents, c.weight_grams, c.in_stock,
    c.variety, c.species,
    r.name AS roaster_name, r.slug AS roaster_slug,
    co.name AS country_name,
    reg.name AS region_name, reg.id AS coffee_region_id,
    p.name AS producer_name, p.id AS coffee_producer_id
FROM coffees c
JOIN roasters r ON r.id = c.roaster_id
LEFT JOIN countries co ON co.code = c.country_code
LEFT JOIN regions reg ON reg.id = c.region_id
LEFT JOIN producers p ON p.id = c.producer_id
WHERE c.region_id = $1 AND r.opted_out = false
ORDER BY c.name;

-- name: ListCoffeesByProducer :many
SELECT
    c.id, c.roaster_id, c.name, c.product_url, c.image_url,
    c.country_code, c.process, c.roast_level,
    c.tasting_notes, c.price_cents, c.weight_grams, c.in_stock,
    c.variety, c.species,
    r.name AS roaster_name, r.slug AS roaster_slug,
    co.name AS country_name,
    reg.name AS region_name, reg.id AS coffee_region_id,
    p.name AS producer_name, p.id AS coffee_producer_id
FROM coffees c
JOIN roasters r ON r.id = c.roaster_id
LEFT JOIN countries co ON co.code = c.country_code
LEFT JOIN regions reg ON reg.id = c.region_id
LEFT JOIN producers p ON p.id = c.producer_id
WHERE c.producer_id = $1 AND r.opted_out = false
ORDER BY c.name;

-- name: ListCoffeesNeedingBackfill :many
SELECT c.id, c.origin_raw, c.region_raw
FROM coffees c
WHERE c.country_code IS NULL AND c.origin_raw IS NOT NULL AND c.origin_raw != '';

-- name: UpdateCoffeeOrigin :exec
UPDATE coffees
SET country_code = $2, region_id = $3
WHERE id = $1;

-- name: ListCoffeesNeedingVariety :many
SELECT c.id, c.variety_raw
FROM coffees c
WHERE c.variety_raw IS NOT NULL AND c.variety_raw != '' AND c.variety IS NULL;

-- name: UpdateCoffeeVariety :exec
UPDATE coffees
SET variety = $2, species = $3
WHERE id = $1;

-- name: ListCoffeesForSimilarity :many
SELECT
    c.id, c.tasting_notes, c.process, c.roast_level, c.variety,
    c.region_id, reg.latitude, reg.longitude
FROM coffees c
JOIN roasters r ON r.id = c.roaster_id
LEFT JOIN regions reg ON reg.id = c.region_id
WHERE c.in_stock = true AND r.opted_out = false;
