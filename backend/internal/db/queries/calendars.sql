-- name: CreateCalendar :one
INSERT INTO calendars (id, title, description, dates)
VALUES ($1, $2, $3, $4)
RETURNING id, title, description, dates, created_at;

-- name: GetCalendarByID :one
SELECT id, title, description, dates, created_at
FROM calendars
WHERE id = $1;

-- name: ListCalendars :many
SELECT id, title, description, dates, created_at
FROM calendars
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

