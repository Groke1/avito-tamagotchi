-- +goose Up
CREATE SCHEMA IF NOT EXISTS users;

CREATE TABLE users.user_streaks (
    user_id UUID PRIMARY KEY
        REFERENCES users.users(id) ON DELETE CASCADE,
    current_streak INTEGER NOT NULL DEFAULT 0 CHECK (current_streak >= 0),
    last_active_date DATE NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION users.update_user_streaks_timestamp() RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = now();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE OR REPLACE TRIGGER trigger_update_user_streaks_timestamp
BEFORE UPDATE ON users.user_streaks
FOR EACH ROW
EXECUTE FUNCTION users.update_user_streaks_timestamp();

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trigger_update_user_streaks_timestamp ON users.user_streaks;
DROP FUNCTION IF EXISTS users.update_user_streaks_timestamp();
DROP TABLE IF EXISTS users.user_streaks;
-- +goose StatementEnd
