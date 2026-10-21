-- +goose Up
UPDATE posts SET description = '' WHERE description IS NULL;

ALTER TABLE posts
ALTER COLUMN description SET NOT NULL;

ALTER TABLE posts
ALTER COLUMN published_at DROP NOT NULL;

-- +goose Down
ALTER TABLE posts
ALTER COLUMN description DROP NOT NULL;

UPDATE posts SET published_at = created_at WHERE published_at IS NULL;

ALTER TABLE posts
ALTER COLUMN published_at SET NOT NULL;