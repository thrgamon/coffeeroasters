-- name: AddUserTastingNote :one
INSERT INTO user_tasting_notes (user_id, coffee_id, tasting_note)
VALUES ($1, $2, $3)
ON CONFLICT (user_id, coffee_id, tasting_note) DO NOTHING
RETURNING *;

-- name: RemoveUserTastingNote :exec
DELETE FROM user_tasting_notes
WHERE user_id = $1 AND coffee_id = $2 AND tasting_note = $3;

-- name: ListUserTastingNotesForCoffee :many
SELECT tasting_note FROM user_tasting_notes
WHERE user_id = $1 AND coffee_id = $2;

-- name: ListCrowdsourcedTastingNotes :many
SELECT tasting_note, COUNT(*) AS vote_count
FROM user_tasting_notes
WHERE coffee_id = $1
GROUP BY tasting_note
ORDER BY vote_count DESC, tasting_note ASC;
