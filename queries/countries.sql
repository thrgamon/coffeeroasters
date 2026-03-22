-- name: GetCountryByCode :one
SELECT code, name, alpha3
FROM countries
WHERE code = $1;

-- name: ListCountriesWithCoffeeCount :many
SELECT co.code, co.name, count(DISTINCT ac.id)::int AS coffee_count
FROM countries co
JOIN (
    SELECT c.id, c.country_code AS match_code, c.roaster_id
    FROM coffees c
    WHERE c.in_stock = true
    UNION
    SELECT bc_c.id, bc.country_code AS match_code, bc_c.roaster_id
    FROM blend_components bc
    JOIN coffees bc_c ON bc_c.id = bc.coffee_id
        AND bc_c.in_stock = true AND bc_c.is_blend = true
) ac ON ac.match_code = co.code
JOIN roasters r ON r.id = ac.roaster_id AND r.opted_out = false
GROUP BY co.code, co.name
ORDER BY co.name;

-- name: ListRegionsByCountry :many
SELECT reg.id, reg.name, count(c.id)::int AS coffee_count
FROM regions reg
LEFT JOIN coffees c ON c.region_id = reg.id AND c.in_stock = true
LEFT JOIN roasters r ON r.id = c.roaster_id AND r.opted_out = false
WHERE reg.country_code = $1
GROUP BY reg.id, reg.name
ORDER BY reg.name;

-- name: GetRegionByID :one
SELECT reg.id, reg.country_code, reg.name, co.name AS country_name, reg.latitude, reg.longitude
FROM regions reg
JOIN countries co ON co.code = reg.country_code
WHERE reg.id = $1;

-- name: GetOrCreateRegion :one
INSERT INTO regions (country_code, name)
VALUES ($1, $2)
ON CONFLICT (country_code, name) DO UPDATE SET name = EXCLUDED.name
RETURNING id;

-- name: UpdateRegionCoordinates :exec
UPDATE regions SET latitude = $2, longitude = $3 WHERE id = $1;

-- name: ListRegionsNeedingGeocode :many
SELECT r.id, r.name, c.name AS country_name
FROM regions r
JOIN countries c ON c.code = r.country_code
WHERE r.latitude IS NULL;

-- name: ListNearbyRegions :many
SELECT
    r.id,
    r.name,
    r.country_code,
    c.name AS country_name,
    (6371 * acos(
        cos(radians(@source_lat::float8)) * cos(radians(r.latitude)) *
        cos(radians(r.longitude) - radians(@source_lon::float8)) +
        sin(radians(@source_lat::float8)) * sin(radians(r.latitude))
    ))::int AS distance_km,
    count(cof.id)::int AS coffee_count
FROM regions r
JOIN countries c ON c.code = r.country_code
LEFT JOIN coffees cof ON cof.region_id = r.id
LEFT JOIN roasters ro ON ro.id = cof.roaster_id AND ro.opted_out = false
WHERE r.id != @exclude_region_id
  AND r.latitude IS NOT NULL
  AND r.longitude IS NOT NULL
  AND (6371 * acos(
        cos(radians(@source_lat::float8)) * cos(radians(r.latitude)) *
        cos(radians(r.longitude) - radians(@source_lon::float8)) +
        sin(radians(@source_lat::float8)) * sin(radians(r.latitude))
    )) <= @max_distance_km::int
GROUP BY r.id, r.name, r.country_code, c.name
ORDER BY distance_km;
