-- name: GetCountryByCode :one
SELECT code, name, alpha3
FROM countries
WHERE code = $1;

-- name: ListCountriesWithCoffeeCount :many
SELECT co.code, co.name, count(c.id)::int AS coffee_count
FROM countries co
JOIN coffees c ON c.country_code = co.code
JOIN roasters r ON r.id = c.roaster_id
WHERE r.opted_out = false
GROUP BY co.code, co.name
ORDER BY co.name;

-- name: ListRegionsByCountry :many
SELECT reg.id, reg.name, count(c.id)::int AS coffee_count
FROM regions reg
JOIN coffees c ON c.region_id = reg.id
JOIN roasters r ON r.id = c.roaster_id
WHERE reg.country_code = $1 AND r.opted_out = false
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
