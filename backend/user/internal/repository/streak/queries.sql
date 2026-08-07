-- name: GetStreakByUserIDForUpdate :one
SELECT user_id, current_streak, last_active_date
FROM users.user_streaks
WHERE user_id = sqlc.arg(user_id)
FOR UPDATE;

-- name: UpdateStreak :exec
INSERT INTO users.user_streaks (
    user_id, current_streak, last_active_date
)
VALUES (
sqlc.arg(user_id),
   sqlc.arg(current_streak),
   sqlc.arg(last_active_date)
)
ON CONFLICT (user_id) DO UPDATE
    SET current_streak = sqlc.arg(current_streak),
    last_active_date = sqlc.arg(last_active_date)
    RETURNING current_streak, last_active_date;