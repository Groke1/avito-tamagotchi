-- name: GetTaskByID :one
SELECT id, title, description, reward_coins, reward_xp
FROM tasks
WHERE id = $1;

-- name: FindByIDForUser :one
SELECT
    t.id,
    t.title,
    t.description,
    t.reward_coins,
    t.reward_xp,
    COALESCE(ut.status, 'active')::text AS status,
    ut.completed_at
FROM tasks t
LEFT JOIN user_tasks ut ON ut.task_id = t.id AND ut.user_id = $1
WHERE t.id = $2;

-- name: ListRandomTasksForUser :many
SELECT
    t.id,
    t.title,
    t.description,
    t.reward_coins,
    t.reward_xp,
    'active' AS status
FROM tasks t
ORDER BY RANDOM()
LIMIT $1;

-- name: CreateUserTasksBatch :exec
INSERT INTO user_tasks (user_id, task_id)
SELECT 
    sqlc.arg(user_id)::uuid AS user_id, 
    unnest(sqlc.arg(task_ids)::uuid[]) AS task_id
;
-- RETURNING user_id, task_id, updated_at;

-- name: GetUserTaskForUpdate :one
SELECT status, completed_at
FROM user_tasks
WHERE user_id = $1 AND task_id = $2
FOR UPDATE;

-- name: InsertUserTaskCompleted :exec
INSERT INTO user_tasks (user_id, task_id, status, completed_at, updated_at)
VALUES ($1, $2, 'completed', $3, NOW());

-- name: UpdateUserTaskCompleted :exec
UPDATE user_tasks
SET status = 'completed', completed_at = $1, updated_at = NOW()
WHERE user_id = $2 AND task_id = $3;