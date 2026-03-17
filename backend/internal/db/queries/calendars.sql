-- name: CreateCalendar :one
INSERT INTO calendars (
  title,
  description,
  location,
  accept_responses_until,
  password,
  salt,
  admin_token
)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, admin_token;

-- name: GetCalendarByID :one
SELECT 
  id, 
  title,
  description,
  location,
  accept_responses_until,
  created_at,
  updated_at
FROM calendars
WHERE id = $1;

-- name: GetCalendarByIDAndAdminToken :one
SELECT 
  id, 
  title,
  description,
  location,
  accept_responses_until,
  created_at,
  updated_at
FROM calendars
WHERE id = $1 AND admin_token = $2;

-- name: DeleteCalendarByIDAndAdminToken :exec
DELETE FROM calendars
WHERE id = $1 AND admin_token = $2;