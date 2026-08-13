CREATE TYPE trip_status AS ENUM (
    'in_progress',
    'pending_delivery',
    'delivered'
);

CREATE TABLE IF NOT EXISTS trip_events (
    id SERIAL PRIMARY KEY,
    description TEXT NOT NULL,
    is_negative BOOLEAN NOT NULL
);

CREATE TABLE IF NOT EXISTS pet_trips (
    id BIGSERIAL PRIMARY KEY,
    pet_id BIGINT NOT NULL REFERENCES pets(id),
    user_id UUID NOT NULL,
    location VARCHAR(100) NOT NULL DEFAULT 'Zootopia',
    reward_xp INTEGER,
    reward_coins INTEGER,
    reward_code TEXT,
    story TEXT NOT NULL,
    status trip_status NOT NULL DEFAULT 'in_progress',
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ended_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_pet_trips_pet_id_started_at
    ON pet_trips (pet_id, started_at DESC);

CREATE INDEX idx_pet_trips_finished ON pet_trips (ended_at)
    WHERE status = 'in_progress';
;