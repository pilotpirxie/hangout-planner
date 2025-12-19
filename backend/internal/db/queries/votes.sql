-- name: CreateVote :one
INSERT INTO votes (id, calendar_id, username, available)
VALUES ($1, $2, $3, $4)
RETURNING id, calendar_id, username, available, created_at;

-- name: ListVotesByCalendar :many
SELECT id, calendar_id, username, available, created_at
FROM votes
WHERE calendar_id = $1
ORDER BY created_at ASC;

