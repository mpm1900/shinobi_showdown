-- +goose up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE teams (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    team_config jsonb,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- +goose down
DROP TABLE teams;
DROP EXTENSION IF EXISTS "pgcrypto";
