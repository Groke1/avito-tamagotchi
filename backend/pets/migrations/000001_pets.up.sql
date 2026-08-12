CREATE TABLE IF NOT EXISTS pets (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE,
    name VARCHAR(32) NOT NULL,
    level INT NOT NULL DEFAULT 1,
    xp INT NOT NULL DEFAULT 0,
    next_level_xp INT NOT NULL DEFAULT 100,
    satiety INT NOT NULL DEFAULT 60,
    happiness INT NOT NULL DEFAULT 75,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_calculated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_feed_at TIMESTAMPTZ,
    last_stroke_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_pets_leaderboard 
ON pets (level DESC, xp DESC);

CREATE TABLE IF NOT EXISTS pets_daily_xp (
    pet_id BIGINT NOT NULL REFERENCES pets(id) ON DELETE CASCADE,
    date DATE NOT NULL DEFAULT CURRENT_DATE,
    gained_xp INT NOT NULL DEFAULT 0,
    PRIMARY KEY (pet_id, date)
);

CREATE INDEX IF NOT EXISTS idx_pets_daily_xp_leaderboard 
ON pets_daily_xp (date, gained_xp DESC);

CREATE TABLE pet_trips (
    id BIGSERIAL PRIMARY KEY,
    pet_id BIGSERIAL NOT NULL REFERENCES pets(id),

    location VARCHAR(100) NOT NULL DEFAULT 'Zootopia',

    started_at TIMESTAMP NOT NULL,
    ended_at TIMESTAMP NOT NULL,

    status VARCHAR(20) NOT NULL,

    reward_xp INT 700,
    reward_coins INT 50,

    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_pet_trips_pet_id_started_at
    ON pet_trips (pet_id, started_at DESC);