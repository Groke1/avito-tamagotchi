-- name: GetUserEventsAndMarkDelivered :many
UPDATE users.user_events AS ue
SET delivered_at = NOW()
WHERE ue.id IN (
    SELECT e.id FROM users.user_events AS e
    WHERE e.user_id = sqlc.arg(user_id) AND e.delivered_at IS NULL
    ORDER BY e.created_at ASC
)
    RETURNING ue.id, ue.user_id, ue.type, ue.xp,
    ue.coins, ue.streak, ue.user_reward_id, ue.created_at;

-- name: AddUserEvent :exec
INSERT INTO users.user_events (
    user_id, type,
    xp, coins,
    streak, user_reward_id)
VALUES (
           sqlc.arg(user_id),sqlc.arg(type),
           sqlc.arg(xp),sqlc.arg(coins),
           sqlc.narg(streak), sqlc.narg(user_reward_id)
);

-- name: DeleteDeliveredEvents :exec
DELETE FROM users.user_events
       WHERE delivered_at IS NOT NULL
        AND delivered_at < NOW() - INTERVAL '30 days';

