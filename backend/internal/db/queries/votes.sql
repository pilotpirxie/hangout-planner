-- name: CreateVote :one
INSERT INTO votes (
  calendar_id,
  calendar_time_slot_id,
  username
)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetVotesByCalendarID :many
SELECT 
  votes.id, 
  votes.calendar_id, 
  votes.calendar_time_slot_id, 
  votes.username, 
  votes.created_at, 
  votes.updated_at, 
  calendar_time_slots.start_date, 
  calendar_time_slots.end_date
FROM votes
LEFT JOIN calendar_time_slots ON votes.calendar_time_slot_id = calendar_time_slots.id
WHERE votes.calendar_id = $1
ORDER BY votes.created_at ASC;

-- name: DeleteVotesByID :exec
DELETE FROM votes
WHERE id = $1;