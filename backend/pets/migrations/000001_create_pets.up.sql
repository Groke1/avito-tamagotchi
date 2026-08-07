CREATE TABLE IF NOT EXISTS pets (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE,
    name VARCHAR(32) NOT NULL,
    level INT NOT NULL DEFAULT 1,
    xp INT NOT NULL DEFAULT 0,
    next_level_xp INT NOT NULL DEFAULT 100,
    satiety INT NOT NULL DEFAULT 60,
    happiness INT NOT NULL DEFAULT 75,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_pets_user_id 
ON pets(user_id);

CREATE INDEX IF NOT EXISTS idx_pets_leaderboard 
ON pets (level DESC, xp DESC);