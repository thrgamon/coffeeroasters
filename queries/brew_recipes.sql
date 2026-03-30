-- name: CreateBrewRecipe :one
INSERT INTO brew_recipes (
    coffee_id, user_id, title, brew_method,
    dose_grams, water_ml, water_temp_c, grind_size,
    brew_time_seconds, notes, is_public
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: UpdateBrewRecipe :one
UPDATE brew_recipes SET
    title = $3,
    brew_method = $4,
    dose_grams = $5,
    water_ml = $6,
    water_temp_c = $7,
    grind_size = $8,
    brew_time_seconds = $9,
    notes = $10,
    is_public = $11,
    updated_at = now()
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: DeleteBrewRecipe :exec
DELETE FROM brew_recipes
WHERE id = $1 AND user_id = $2;

-- name: GetBrewRecipe :one
SELECT br.*,
       u.email AS user_email
FROM brew_recipes br
JOIN users u ON u.id = br.user_id
WHERE br.id = $1;

-- name: ListBrewRecipesByCoffee :many
SELECT br.*,
       u.email AS user_email
FROM brew_recipes br
JOIN users u ON u.id = br.user_id
WHERE br.coffee_id = $1
    AND (br.is_public = true OR br.user_id = sqlc.narg('viewer_user_id')::int)
ORDER BY br.created_at DESC;

-- name: ListBrewRecipesByUser :many
SELECT br.*,
       c.name AS coffee_name,
       r.name AS roaster_name, r.slug AS roaster_slug
FROM brew_recipes br
JOIN coffees c ON c.id = br.coffee_id
JOIN roasters r ON r.id = c.roaster_id
WHERE br.user_id = $1
ORDER BY br.updated_at DESC;
