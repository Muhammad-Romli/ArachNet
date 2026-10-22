-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetPosts :many
SELECT * FROM posts
ORDER BY created_at DESC
LIMIT $1;

-- name: SearchPosts :many
SELECT *
FROM posts
WHERE (title ILIKE '%' || sqlc.arg(search_text)::text || '%'
       OR url ILIKE '%' || sqlc.arg(search_text)::text || '%')
ORDER BY published_at DESC
LIMIT sqlc.arg(row_limit);


-- name: GetPostByID :one
SELECT * FROM posts WHERE id = $1;
