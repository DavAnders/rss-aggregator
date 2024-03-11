-- name: CreatePost :exec
INSERT INTO posts (created_at, updated_at, title, url, description, published_at, feed_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (url) DO UPDATE 
SET title = EXCLUDED.title,
    description = EXCLUDED.description,
    updated_at = EXCLUDED.updated_at,
    published_at = EXCLUDED.published_at
WHERE posts.updated_at < EXCLUDED.updated_at;

-- name: GetPostsByUsers :many
SELECT p.*
FROM posts p
JOIN feed_follows ff ON p.feed_id = ff.feed_id
WHERE ff.user_id = $1
ORDER BY p.published_at DESC
LIMIT $2;