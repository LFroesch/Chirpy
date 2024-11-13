-- +goose Up
-- Add the new column here using ALTER TABLE
ALTER TABLE users ADD COLUMN hashed_password TEXT NOT NULL DEFAULT 'unset';

-- +goose Down
-- Remove the column here using ALTER TABLE
ALTER TABLE users DROP COLUMN hashed_password;

