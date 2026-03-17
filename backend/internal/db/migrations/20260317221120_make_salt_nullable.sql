-- +goose Up
ALTER TABLE calendars
  ALTER COLUMN salt DROP NOT NULL;

-- +goose Down
ALTER TABLE calendars
  ALTER COLUMN salt SET NOT NULL;