-- +goose Up
CREATE SCHEMA IF NOT EXISTS users;

CREATE TYPE users.reward_status AS ENUM ('active', 'redeemed');

CREATE TABLE IF NOT EXISTS users.reward_definitions (
    id   SERIAL PRIMARY KEY,
    code TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    description TEXT NOT NULL

);

INSERT INTO users.reward_definitions (
    code, name, description
)
VALUES
    (
        'DELIVERY_DISCOUNT_10',
        'Скидка 10% на доставку',
        'Скидка 10% на одну покупку с Авито Доставкой'
    ),
    (
        'FREE_LISTING_PROMOTION',
        'Бесплатное продвижение',
        'Одно бесплатное продвижение выбранного объявления'
    ),
    (
        'AUTOTEKA_DISCOUNT_20',
        'Скидка 20% на Автотеку',
        'Скидка 20% на один отчёт об истории автомобиля'
    ),
    (
        'FREE_LISTING_HIGHLIGHT',
        'Выделение объявления',
        'Бесплатное визуальное выделение одного объявления на ограниченный срок'
    ),
    (
        'LISTING_DISCOUNT_15',
        'Скидка 15% на размещение',
        'Скидка 15% на одно платное размещение объявления'
    );

CREATE TABLE IF NOT EXISTS users.user_rewards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users.users(id)
             ON DELETE CASCADE,
    reward_id INTEGER NOT NULL REFERENCES users.reward_definitions(id),
    promo_code TEXT NOT NULL UNIQUE,
    status users.reward_status NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    redeemed_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ
);

CREATE INDEX user_rewards_user_id_idx
    ON users.user_rewards (user_id);

-- +goose Down
DROP INDEX IF EXISTS user_rewards_user_id_idx;
DROP TABLE IF EXISTS users.user_rewards;
DROP TABLE IF EXISTS users.reward_definitions;
DROP TYPE IF EXISTS users.reward_status;
