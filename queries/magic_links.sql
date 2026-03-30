-- name: CreateMagicLink :one
INSERT INTO magic_links (email, token, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetMagicLinkByToken :one
SELECT * FROM magic_links
WHERE token = $1 AND expires_at > now() AND used_at IS NULL;

-- name: MarkMagicLinkUsed :exec
UPDATE magic_links SET used_at = now()
WHERE token = $1;

-- name: DeleteExpiredMagicLinks :exec
DELETE FROM magic_links
WHERE expires_at <= now() OR used_at IS NOT NULL;
