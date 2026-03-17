-- name: InsertBlendComponent :exec
INSERT INTO blend_components (coffee_id, country_code, region_id, variety, percentage)
VALUES ($1, $2, $3, $4, $5);

-- name: DeleteBlendComponents :exec
DELETE FROM blend_components WHERE coffee_id = $1;

-- name: ListBlendComponents :many
SELECT
    bc.id, bc.coffee_id, bc.country_code, bc.region_id,
    bc.variety, bc.percentage,
    co.name AS country_name,
    reg.name AS region_name
FROM blend_components bc
LEFT JOIN countries co ON co.code = bc.country_code
LEFT JOIN regions reg ON reg.id = bc.region_id
WHERE bc.coffee_id = $1
ORDER BY bc.id;
