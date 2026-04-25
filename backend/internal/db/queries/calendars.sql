-- name: CreateCalendar :one
INSERT INTO calendars (
  title,
  description,
  password,
  salt,
  admin_token
)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, admin_token;

-- name: GetCalendarByID :one
SELECT 
  id, 
  title,
  description,
  created_at,
  updated_at
FROM calendars
WHERE id = $1;

-- name: GetCalendarByIDAndAdminToken :one
SELECT 
  id, 
  title,
  description,
  created_at,
  updated_at
FROM calendars
WHERE id = $1 AND admin_token = $2;

-- name: DeleteCalendarByIDAndAdminToken :exec
DELETE FROM calendars
WHERE id = $1 AND admin_token = $2;

-- name: GetCalendarPasswordAndSaltByID :one
SELECT 
  password,
  salt
FROM calendars
WHERE id = $1;