-- name: GetCalendarTimeSlotsByCalendarID :many
SELECT 
  id,
  calendar_id,
  start_date,
  end_date,
  created_at,
  updated_at
FROM calendar_time_slots
WHERE calendar_id = $1
ORDER BY start_date, end_date;

-- name: GetCalendarTimeSlotByID :one
SELECT 
  id,
  calendar_id,
  start_date,
  end_date,
  created_at,
  updated_at
FROM calendar_time_slots
WHERE id = $1;

-- name: CreateCalendarTimeSlot :one
INSERT INTO calendar_time_slots (
  calendar_id,
  start_date,
  end_date
)
VALUES ($1, $2, $3)
RETURNING *;

-- name: DeleteCalendarTimeSlotByID :exec
DELETE FROM calendar_time_slots
WHERE id = $1;
