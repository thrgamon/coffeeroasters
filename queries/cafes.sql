-- name: UpsertCafe :one
INSERT INTO cafes (roaster_id, slug, name, address, suburb, state, postcode, latitude, longitude, phone, instagram, website_url, image_url)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
ON CONFLICT (roaster_id, slug) DO UPDATE SET
    name = EXCLUDED.name,
    address = EXCLUDED.address,
    suburb = EXCLUDED.suburb,
    state = EXCLUDED.state,
    postcode = EXCLUDED.postcode,
    latitude = EXCLUDED.latitude,
    longitude = EXCLUDED.longitude,
    phone = EXCLUDED.phone,
    instagram = EXCLUDED.instagram,
    website_url = EXCLUDED.website_url,
    image_url = EXCLUDED.image_url,
    updated_at = now()
RETURNING id;

-- name: GetCafeBySlug :one
SELECT c.id, c.roaster_id, c.slug, c.name, c.address, c.suburb, c.state, c.postcode,
       c.latitude, c.longitude, c.phone, c.instagram, c.website_url, c.image_url,
       r.name AS roaster_name, r.slug AS roaster_slug
FROM cafes c
JOIN roasters r ON r.id = c.roaster_id
WHERE c.slug = $1 AND c.active = true;

-- name: ListCafes :many
SELECT c.id, c.slug, c.name, c.address, c.suburb, c.state, c.postcode,
       c.latitude, c.longitude, c.phone, c.instagram, c.website_url, c.image_url,
       r.id AS roaster_id, r.name AS roaster_name, r.slug AS roaster_slug
FROM cafes c
JOIN roasters r ON r.id = c.roaster_id
WHERE c.active = true AND r.active = true AND r.opted_out = false
ORDER BY c.state, c.name;

-- name: ListCafesByState :many
SELECT c.id, c.slug, c.name, c.address, c.suburb, c.state, c.postcode,
       c.latitude, c.longitude, c.phone, c.instagram, c.website_url, c.image_url,
       r.id AS roaster_id, r.name AS roaster_name, r.slug AS roaster_slug
FROM cafes c
JOIN roasters r ON r.id = c.roaster_id
WHERE c.state = $1 AND c.active = true AND r.active = true AND r.opted_out = false
ORDER BY c.name;

-- name: ListCafesByRoaster :many
SELECT c.id, c.slug, c.name, c.address, c.suburb, c.state, c.postcode,
       c.latitude, c.longitude, c.phone, c.instagram, c.website_url, c.image_url
FROM cafes c
WHERE c.roaster_id = $1 AND c.active = true
ORDER BY c.name;

-- name: CountCafes :one
SELECT count(*) FROM cafes WHERE active = true;
