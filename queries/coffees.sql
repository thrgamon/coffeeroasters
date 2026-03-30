-- name: UpsertCoffee :one
-- Returns true if the row was newly inserted, false if updated.
INSERT INTO coffees (
    roaster_id, name, product_url, image_url,
    origin_raw, region_raw, variety_raw, process_raw, roast_raw,
    tasting_notes_raw, price_raw, weight_raw, currency, in_stock,
    process, roast_level, tasting_notes, price_cents, weight_grams,
    country_code, region_id, producer_id, producer_raw,
    variety, species,
    price_per_100g_min, price_per_100g_max, is_blend,
    description, brew_recipe_raw, source_hash, is_decaf,
    last_seen_at
)
VALUES (
    $1, $2, $3, $4,
    $5, $6, $7, $8, $9,
    $10, $11, $12, $13, $14,
    $15, $16, $17, $18, $19,
    $20, $21, $22, $23,
    $24, $25,
    $26, $27, $28,
    $29, $30, $31, $32,
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
    price_per_100g_min = EXCLUDED.price_per_100g_min,
    price_per_100g_max = EXCLUDED.price_per_100g_max,
    is_blend = EXCLUDED.is_blend,
    description = EXCLUDED.description,
    brew_recipe_raw = EXCLUDED.brew_recipe_raw,
    source_hash = EXCLUDED.source_hash,
    is_decaf = EXCLUDED.is_decaf,
    last_seen_at = now(),
    last_changed_at = CASE
        WHEN coffees.price_cents IS DISTINCT FROM EXCLUDED.price_cents
            OR coffees.in_stock IS DISTINCT FROM EXCLUDED.in_stock
        THEN now()
        ELSE coffees.last_changed_at
    END,
    updated_at = now()
RETURNING (xmax = 0) AS is_new, id;

-- name: ListCoffeesFiltered :many
SELECT
    c.id, c.roaster_id, c.name, c.product_url, c.image_url,
    c.country_code, c.origin_raw, c.region_raw, c.process, c.roast_level,
    c.tasting_notes, c.price_cents, c.weight_grams, c.in_stock,
    c.variety, c.species,
    c.price_per_100g_min, c.price_per_100g_max, c.is_blend, c.is_decaf,
    r.name AS roaster_name, r.slug AS roaster_slug, r.logo_url AS roaster_logo_url,
    co.name AS country_name,
    reg.name AS region_name, reg.id AS coffee_region_id,
    p.name AS producer_name, p.id AS coffee_producer_id,
    COUNT(*) OVER() AS total_count
FROM coffees c
JOIN roasters r ON r.id = c.roaster_id
LEFT JOIN countries co ON co.code = c.country_code
LEFT JOIN regions reg ON reg.id = c.region_id
LEFT JOIN producers p ON p.id = c.producer_id
WHERE r.opted_out = false
    AND c.in_stock = true
    AND (sqlc.narg('query')::text IS NULL
         OR c.search_vector @@ plainto_tsquery('english', sqlc.narg('query')))
    AND (sqlc.narg('origin')::text IS NULL
         OR c.country_code = sqlc.narg('origin')
         OR (c.is_blend AND EXISTS (
             SELECT 1 FROM blend_components bc
             WHERE bc.coffee_id = c.id AND bc.country_code = sqlc.narg('origin'))))
    AND (sqlc.narg('process')::text IS NULL OR c.process = sqlc.narg('process'))
    AND (sqlc.narg('roast')::text IS NULL OR c.roast_level = sqlc.narg('roast'))
    AND (sqlc.narg('variety')::text IS NULL OR c.variety = sqlc.narg('variety'))
    AND (sqlc.narg('roaster_state')::text IS NULL OR r.state = sqlc.narg('roaster_state'))
    AND (sqlc.narg('decaf')::text IS NULL
         OR (sqlc.narg('decaf') = 'only' AND c.is_decaf = true)
         OR (sqlc.narg('decaf') = 'exclude' AND c.is_decaf = false))
    AND (c.country_code IS NOT NULL OR c.process IS NOT NULL OR array_length(c.tasting_notes, 1) IS NOT NULL)
ORDER BY
    CASE WHEN sqlc.narg('query')::text IS NOT NULL
        THEN ts_rank(c.search_vector, plainto_tsquery('english', sqlc.narg('query')))
    END DESC NULLS LAST,
    c.name
LIMIT sqlc.arg('lim') OFFSET sqlc.arg('off');

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
    c.price_per_100g_min, c.price_per_100g_max, c.is_blend, c.is_decaf,
    r.name AS roaster_name, r.slug AS roaster_slug, r.logo_url AS roaster_logo_url,
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
    c.price_per_100g_min, c.price_per_100g_max, c.is_blend, c.is_decaf,
    c.description, c.brew_recipe_raw,
    c.first_seen_at, c.last_seen_at,
    r.name AS roaster_name, r.slug AS roaster_slug, r.logo_url AS roaster_logo_url,
    co.name AS country_name,
    reg.name AS region_name,
    p.name AS producer_name
FROM coffees c
JOIN roasters r ON r.id = c.roaster_id
LEFT JOIN countries co ON co.code = c.country_code
LEFT JOIN regions reg ON reg.id = c.region_id
LEFT JOIN producers p ON p.id = c.producer_id
WHERE c.id = $1;


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
    c.price_per_100g_min, c.price_per_100g_max, c.is_blend, c.is_decaf,
    r.name AS roaster_name, r.slug AS roaster_slug, r.logo_url AS roaster_logo_url,
    co.name AS country_name,
    reg.name AS region_name, reg.id AS coffee_region_id,
    p.name AS producer_name, p.id AS coffee_producer_id
FROM coffees c
JOIN roasters r ON r.id = c.roaster_id
LEFT JOIN countries co ON co.code = c.country_code
LEFT JOIN regions reg ON reg.id = c.region_id
LEFT JOIN producers p ON p.id = c.producer_id
WHERE r.opted_out = false
    AND (c.country_code = $1
        OR (c.is_blend AND EXISTS (
            SELECT 1 FROM blend_components bc
            WHERE bc.coffee_id = c.id AND bc.country_code = $1)))
ORDER BY c.name;

-- name: ListCoffeesByRegion :many
SELECT
    c.id, c.roaster_id, c.name, c.product_url, c.image_url,
    c.country_code, c.process, c.roast_level,
    c.tasting_notes, c.price_cents, c.weight_grams, c.in_stock,
    c.variety, c.species,
    c.price_per_100g_min, c.price_per_100g_max, c.is_blend, c.is_decaf,
    r.name AS roaster_name, r.slug AS roaster_slug, r.logo_url AS roaster_logo_url,
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
    c.price_per_100g_min, c.price_per_100g_max, c.is_blend, c.is_decaf,
    r.name AS roaster_name, r.slug AS roaster_slug, r.logo_url AS roaster_logo_url,
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
    c.region_id, reg.latitude, reg.longitude,
    c.embedding
FROM coffees c
JOIN roasters r ON r.id = c.roaster_id
LEFT JOIN regions reg ON reg.id = c.region_id
WHERE c.in_stock = true AND r.opted_out = false;

-- name: ListCoffeesNeedingEmbedding :many
SELECT c.id, c.description
FROM coffees c
WHERE c.description IS NOT NULL AND c.description != '' AND c.embedding IS NULL;

-- name: UpdateCoffeeEmbedding :exec
UPDATE coffees SET embedding = $2 WHERE id = $1;

-- name: GetSourceHashesByRoaster :many
SELECT c.product_url, c.source_hash
FROM coffees c
JOIN roasters r ON r.id = c.roaster_id
WHERE r.slug = $1 AND c.source_hash IS NOT NULL AND c.product_url IS NOT NULL;

-- name: UpdateCoffeeSeenAndPrice :exec
UPDATE coffees SET
    last_seen_at = now(),
    in_stock = $3,
    price_raw = $4,
    weight_raw = $5,
    price_cents = $6,
    weight_grams = $7,
    price_per_100g_min = $8,
    price_per_100g_max = $9,
    image_url = $10,
    updated_at = now()
WHERE roaster_id = $1 AND name = $2;

-- name: ListCoffeesForFinder :many
SELECT
    c.id, c.roaster_id, c.name, c.product_url, c.image_url,
    c.country_code, c.process, c.roast_level,
    c.tasting_notes, c.price_cents, c.weight_grams, c.in_stock,
    c.variety, c.species,
    c.price_per_100g_min, c.price_per_100g_max, c.is_blend, c.is_decaf,
    r.name AS roaster_name, r.slug AS roaster_slug, r.logo_url AS roaster_logo_url,
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
    AND (c.country_code IS NOT NULL OR c.process IS NOT NULL OR array_length(c.tasting_notes, 1) IS NOT NULL)
ORDER BY c.name;
