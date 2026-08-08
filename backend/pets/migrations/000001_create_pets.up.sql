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
    last_feed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_stroke_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
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