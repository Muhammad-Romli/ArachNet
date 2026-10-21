-- name: CreateFeed :one
INSERT INTO feeds(id, created_at, updated_at, name, url, user_id) 
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeedByURL :one
SELECT * FROM feeds
WHERE url = $1;

-- name: GetFeedByName :one
SELECT * FROM feeds
WHERE name = $1;

-- name: GetListFeeds :many
SELECT feeds.*, users.name AS user_name
FROM feeds
INNER JOIN users ON feeds.user_id = users.id;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = NOW(),
    last_updated = NOW()
WHERE $1 = name OR $1 = url;