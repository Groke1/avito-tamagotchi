CREATE TABLE IF NOT EXISTS tasks (
    id uuid PRIMARY KEY,
    title text NOT NULL,
    description text NOT NULL,
    reward_coins integer NOT NULL,
    reward_xp bigint NOT NULL
);

CREATE TABLE IF NOT EXISTS user_tasks (
    user_id uuid NOT NULL,
    task_id uuid NOT NULL,
    status text NOT NULL DEFAULT 'active',
    completed_at timestamptz NULL,
    updated_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, task_id)
);

CREATE INDEX IF NOT EXISTS idx_user_tasks_task_id ON user_tasks (task_id);