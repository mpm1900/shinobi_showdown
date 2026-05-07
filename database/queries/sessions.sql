-- name: CreateSession :one
INSERT INTO sessions (user_id, expires_at)
VALUES ($1, $2)
RETURNING *;

-- name: GetSession :one
SELECT * FROM sessions WHERE id = $1;

-- name: DeleteSession :exec
DELETE FROM sessions WHERE id = $1;

-- name: DeleteUserSessions :exec
DELETE FROM sessions WHERE user_id = $1;

-- name: DeleteExpiredSessions :exec
DELETE FROM sessions WHERE expires_at < CURRENT_TIMESTAMP;

-- name: GetUserBySessionID :one
SELECT u.* FROM users u
JOIN sessions s ON s.user_id = u.id
WHERE s.id = $1 AND s.expires_at > CURRENT_TIMESTAMP;
