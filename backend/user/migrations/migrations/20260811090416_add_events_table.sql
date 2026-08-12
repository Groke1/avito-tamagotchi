-- +goose Up

CREATE SCHEMA IF NOT EXISTS users;

CREATE TYPE users.event_type AS ENUM (
    'streak_reward',
    'new_reward'
);

CREATE TABLE users.user_events (
      id BIGSERIAL PRIMARY KEY,
      user_id UUID NOT NULL REFERENCES users.users(id) ON DELETE CASCADE,
      type users.event_type NOT NULL,
      xp INTEGER NOT NULL DEFAULT 0,
      coins INTEGER NOT NULL DEFAULT 0,
      streak INTEGER,
      user_reward_id UUID REFERENCES users.user_rewards(id)
          ON DELETE CASCADE,
      created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
      delivered_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS user_events_user_id_idx
    ON users.user_events (user_id);

-- +goose Down

DROP INDEX IF EXISTS user_events_user_id_idx;
DROP TABLE IF EXISTS users.user_events;
DROP TYPE IF EXISTS users.event_type;

