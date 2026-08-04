-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS account;
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE account.users (
       id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
       username      VARCHAR(100) NOT NULL UNIQUE,
       email         VARCHAR(255) NOT NULL UNIQUE,
       password_hash VARCHAR(255) NOT NULL,
       coins         BIGINT NOT NULL CHECK (coins >= 0),
       created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE account.refresh_tokens (
        id          BIGSERIAL PRIMARY KEY,
        user_id     UUID NOT NULL REFERENCES account.users(id) ON DELETE CASCADE,
        token_hash  VARCHAR(64) NOT NULL UNIQUE,
        expires_at  TIMESTAMPTZ NOT NULL,
        created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_refresh_tokens_user_id ON account.refresh_tokens(user_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_refresh_tokens_user_id;
DROP TABLE IF EXISTS account.refresh_tokens;
DROP TABLE IF EXISTS account.users;

-- +goose StatementEnd
