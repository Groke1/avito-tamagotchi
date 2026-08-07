-- name: AddToken :exec
INSERT INTO users.refresh_tokens (
    user_id,
    token_hash,
    expires_at
)
VALUES (
           sqlc.arg(user_id),
           sqlc.arg(token_hash),
           sqlc.arg(expires_at)
       );

-- name: GetRefreshTokenByHashForUpdate :one
SELECT
    id,
    user_id,
    token_hash,
    expires_at,
    created_at
FROM users.refresh_tokens
WHERE token_hash = sqlc.arg(token_hash)
FOR UPDATE ;

-- name: DeleteRefreshTokenByHash :exec
DELETE FROM users.refresh_tokens
WHERE token_hash = sqlc.arg(token_hash);

-- name: DeleteExpiredTokens :exec
DELETE FROM users.refresh_tokens WHERE expires_at < NOW();

-- name: DeleteSession :exec
DELETE FROM users.refresh_tokens
WHERE user_id = sqlc.arg(user_id)
  AND token_hash = sqlc.arg(token_hash);
