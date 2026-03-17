-- name: InsertScrapeRun :one
INSERT INTO scrape_runs (roaster_id, status)
VALUES ($1, 'running')
RETURNING id;

-- name: CompleteScrapeRun :exec
UPDATE scrape_runs
SET status = $2,
    finished_at = now(),
    coffees_found = $3,
    coffees_added = $4,
    coffees_updated = $5,
    duration_ms = $6
WHERE id = $1;

-- name: FailScrapeRun :exec
UPDATE scrape_runs
SET status = 'failed',
    finished_at = now(),
    error_message = $2,
    duration_ms = $3
WHERE id = $1;

-- name: ListScrapeRunsByRoaster :many
SELECT id, roaster_id, started_at, finished_at, status,
       coffees_found, coffees_added, coffees_updated, coffees_removed,
       pages_visited, error_message, duration_ms
FROM scrape_runs
WHERE roaster_id = $1
ORDER BY started_at DESC
LIMIT $2;
