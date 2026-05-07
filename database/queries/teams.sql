-- name: GetTeamsByUser :many
SELECT
    id,
    user_id,
    COALESCE(team_config, '{}'::jsonb) AS team_config,
    created_at
FROM teams
WHERE user_id = $1;

-- name: GetTeamByID :one
SELECT
    id,
    user_id,
    COALESCE(team_config, '{}'::jsonb) AS team_config,
    created_at
FROM teams
WHERE id = $1;

-- name: CreateTeam :one
INSERT INTO teams (user_id, team_config)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateTeam :one
UPDATE teams
SET team_config = $2
WHERE id = $1
RETURNING *;

-- name: DeleteTeam :exec
DELETE FROM teams WHERE id = $1;
