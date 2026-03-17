-- +goose Up
ALTER TABLE calendars
  ADD COLUMN salt text NOT NULL DEFAULT gen_random_uuid(),
  ADD COLUMN admin_token text NOT NULL DEFAULT gen_random_uuid();

-- +goose Down
ALTER TABLE calendars
  DROP COLUMN salt,
  DROP COLUMN admin_token;