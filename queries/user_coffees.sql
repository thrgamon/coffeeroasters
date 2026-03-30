-- name: UpsertUserCoffee :one
INSERT INTO user_coffees (user_id, coffee_id, status, liked, rating, review, drunk_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (user_id, coffee_id)
DO UPDATE SET status = EXCLUDED.status, liked = EXCLUDED.liked, rating = EXCLUDED.rating,
             review = EXCLUDED.review, drunk_at = EXCLUDED.drunk_at, updated_at = now()
RETURNING *;

-- name: DeleteUserCoffee :exec
DELETE FROM user_coffees
WHERE user_id = $1 AND coffee_id = $2;

-- name: ListUserCoffees :many
SELECT uc.*, c.name AS coffee_name, c.image_url AS coffee_image_url,
       r.name AS roaster_name, r.slug AS roaster_slug, r.logo_url AS roaster_logo_url,
       c.process, c.roast_level, c.tasting_notes, c.variety,
       c.price_cents, c.weight_grams, c.price_per_100g_min, c.price_per_100g_max,
       c.in_stock, c.is_blend, c.is_decaf, c.product_url,
       c.country_code, co.name AS country_name,
       c.region_id, reg.name AS region_name,
       c.species, c.description
FROM user_coffees uc
JOIN coffees c ON c.id = uc.coffee_id
JOIN roasters r ON r.id = c.roaster_id
LEFT JOIN countries co ON co.code = c.country_code
LEFT JOIN regions reg ON reg.id = c.region_id
WHERE uc.user_id = $1
ORDER BY uc.updated_at DESC;

-- name: GetUserCoffee :one
SELECT * FROM user_coffees
WHERE user_id = $1 AND coffee_id = $2;

-- name: ListUserCoffeeIDs :many
SELECT coffee_id, status, liked, rating FROM user_coffees
WHERE user_id = $1;

-- name: CreateUserPasswordless :one
INSERT INTO users (email, password_hash)
VALUES ($1, '')
RETURNING *;
