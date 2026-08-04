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