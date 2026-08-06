-- name: AddUser :one
INSERT INTO users.users (
    username,
    email,
    password_hash,
    coins
)
VALUES (
sqlc.arg(username),
sqlc.arg(email),
sqlc.arg(password_hash),
sqlc.arg(coins)
) RETURNING id;

-- name: GetUserByID :one
SELECT id, username, email, coins
FROM users.users WHERE id = sqlc.arg(id);

-- name: GetUserByEmail :one
SELECT
    id,
    username,
    email,
    password_hash,
    coins
FROM users.users
WHERE email = sqlc.arg(email);

-- name: UpdateCoins :one
UPDATE users.users
SET coins = coins + sqlc.arg(delta_coins)
WHERE id = sqlc.arg(user_id)
  AND coins + sqlc.arg(delta_coins) >= 0
RETURNING coins;

-- name: GetUsersByIDs :many
SELECT id, username
FROM users.users
WHERE id = ANY(sqlc.arg(ids)::uuid[]);

-- name: UserExists :one
SELECT EXISTS (
    SELECT 1
    FROM users.users
    WHERE id = sqlc.arg(user_id)
);
