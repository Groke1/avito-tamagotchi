-- name: AddUser :one
INSERT INTO account.users (
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
FROM account.users WHERE id = sqlc.arg(id);

-- name: GetUserByEmail :one
SELECT
    id,
    username,
    email,
    password_hash,
    coins
FROM account.users
WHERE email = sqlc.arg(email);

-- name: UpdateCoins :one
UPDATE account.users
SET coins = coins + sqlc.arg(delta_coins)
WHERE id = sqlc.arg(user_id)
  AND coins + sqlc.arg(delta_coins) >= 0
RETURNING coins;

-- name: GetUsersByIDs :many
SELECT id, username
FROM account.users
WHERE id = ANY(sqlc.arg(ids)::uuid[]);

-- name: UserExists :one
SELECT EXISTS (
    SELECT 1
    FROM account.users
    WHERE id = sqlc.arg(user_id)
);
