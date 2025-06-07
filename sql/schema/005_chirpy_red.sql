-- +goose Up
-- Add the new column here using ALTER TABLE
ALTER TABLE users
ADD COLUMN is_chirpy_red BOOLEAN NOT NULL
DEFAULT FALSE;

-- +goose Down
-- Remove the column here using ALTER TABLE
ALTER TABLE users DROP COLUMN is_chirpy_red;