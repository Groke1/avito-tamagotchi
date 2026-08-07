-- name: GetUserRewardsByUserID :many
SELECT
    ur.id, ur.user_id, ur.promo_code,
    rd.name, rd.description,
    ur.status, ur.redeemed_at,
    ur.expires_at
FROM users.user_rewards AS ur
         JOIN users.reward_definitions AS rd
              ON rd.id = ur.reward_id
WHERE ur.user_id = sqlc.arg(user_id)
ORDER BY ur.created_at DESC;

-- name: GetActiveRewardsByUserID :many
SELECT ur.id, ur.user_id, ur.promo_code,
       rd.name, rd.description,
       ur.status, ur.redeemed_at,
       ur.expires_at
FROM users.user_rewards ur JOIN users.reward_definitions rd
                                ON ur.reward_id = rd.id
WHERE user_id = sqlc.arg(user_id)
  AND status = 'active'
  AND (
        ur.expires_at IS NULL OR ur.expires_at > NOW()
    )
ORDER BY ur.created_at DESC;;

-- name: GetRewardByUserIDAndRewardID :one
SELECT ur.id, ur.user_id, ur.promo_code,
       rd.name, rd.description,
       ur.status, ur.redeemed_at,
       ur.expires_at
FROM users.user_rewards ur JOIN users.reward_definitions rd
                                ON ur.reward_id = rd.id
WHERE ur.id = sqlc.arg(reward_id) AND ur.user_id = sqlc.arg(user_id);

-- name: GetRewardDefinitionByCode :one
SELECT
    id, code, name, description
FROM users.reward_definitions
WHERE code = sqlc.arg(code);

-- name: AddUserReward :one
INSERT INTO users.user_rewards (
    user_id, reward_id,
    promo_code, expires_at
)
VALUES (
       sqlc.arg(user_id),
       sqlc.arg(reward_id),
       sqlc.arg(promo_code),
       sqlc.narg(expires_at)
       )
RETURNING id, user_id, promo_code, status, redeemed_at, expires_at;

-- name: RedeemUserReward :one
UPDATE users.user_rewards
SET
    status = 'redeemed',
    redeemed_at = NOW()
WHERE promo_code = sqlc.arg(promo_code)
    AND user_id = sqlc.arg(user_id)
    AND status = 'active'
    AND (
        expires_at IS NULL
        OR expires_at > NOW()
    )
RETURNING id;

