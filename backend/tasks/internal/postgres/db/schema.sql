CREATE TABLE IF NOT EXISTS tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    task_type VARCHAR(20),
    reward_coins INT NOT NULL,
    reward_xp BIGINT NOT NULL,
    description TEXT NOT NULL,
    finished_desc TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS user_tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    task_id UUID NOT NULL REFERENCES tasks(id),
    status TEXT NOT NULL DEFAULT 'active',
    completed_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);