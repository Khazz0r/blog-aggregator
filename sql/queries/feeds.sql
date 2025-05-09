-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: GetFeedByUrl :one
SELECT * FROM feeds
WHERE feeds.url = $1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET
    last_fetched_at = $1,
    updated_at = $2
WHERE feeds.id = $3;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds
ORDER BY feeds.last_fetched_at NULLS FIRST
LIMIT 1;
