-- +goose Up
ALTER TABLE calendars
  DROP COLUMN location,
  DROP COLUMN accept_responses_until;

-- +goose Down
ALTER TABLE calendars
  ADD COLUMN location text,
  ADD COLUMN accept_responses_until timestamptz;