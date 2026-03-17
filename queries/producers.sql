-- name: GetOrCreateProducer :one
INSERT INTO producers (name, country_code, region_id)
VALUES ($1, $2, $3)
ON CONFLICT (name, country_code) DO UPDATE SET
    region_id = COALESCE(EXCLUDED.region_id, producers.region_id)
RETURNING id;

-- name: GetProducerByID :one
SELECT p.id, p.name, p.country_code, p.region_id,
       co.name AS country_name,
       reg.name AS region_name
FROM producers p
LEFT JOIN countries co ON co.code = p.country_code
LEFT JOIN regions reg ON reg.id = p.region_id
WHERE p.id = $1;

-- name: ListProducersWithCoffeeCount :many
SELECT p.id, p.name, p.country_code, co.name AS country_name,
       count(c.id)::int AS coffee_count
FROM producers p
LEFT JOIN countries co ON co.code = p.country_code
JOIN coffees c ON c.producer_id = p.id
JOIN roasters r ON r.id = c.roaster_id
WHERE r.opted_out = false
GROUP BY p.id, p.name, p.country_code, co.name
ORDER BY p.name;
