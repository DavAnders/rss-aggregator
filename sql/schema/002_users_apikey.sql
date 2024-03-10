-- +goose Up
ALTER TABLE users ADD COLUMN api_key VARCHAR(64) NOT NULL UNIQUE;

-- +goose Down
ALTER TABLE users DROP COLUMN api_key;