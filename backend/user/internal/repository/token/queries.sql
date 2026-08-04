-- name: AddToken :exec
INSERT INTO account.refresh_tokens (
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
FROM account.refresh_tokens
WHERE token_hash = sqlc.arg(token_hash)
FOR UPDATE ;

-- name: DeleteRefreshTokenByHash :exec
DELETE FROM account.refresh_tokens
WHERE token_hash = sqlc.arg(token_hash);

-- name: DeleteExpiredTokens :exec
DELETE FROM account.refresh_tokens WHERE expires_at < NOW();
