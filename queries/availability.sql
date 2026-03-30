-- name: UpsertAvailabilityLog :exec
INSERT INTO coffee_availability_log (coffee_id, in_stock, price_cents, recorded_at)
VALUES ($1, $2, $3, CURRENT_DATE)
ON CONFLICT (coffee_id, recorded_at) DO UPDATE SET
    in_stock = EXCLUDED.in_stock,
    price_cents = EXCLUDED.price_cents;

-- name: ListAvailabilityHistory :many
SELECT coffee_id, in_stock, price_cents, recorded_at
FROM coffee_availability_log
WHERE coffee_id = $1
ORDER BY recorded_at DESC
LIMIT $2;

-- name: MarkStaleCoffeesOutOfStock :exec
-- Coffees not seen by the scraper in over 3 days are marked out of stock.
UPDATE coffees
SET in_stock = false, updated_at = now()
WHERE in_stock = true
  AND last_seen_at < now() - interval '3 days';

-- name: ListCoffeesByRoasterWithAvailability :many
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
ORDER BY c.in_stock DESC, c.name;
