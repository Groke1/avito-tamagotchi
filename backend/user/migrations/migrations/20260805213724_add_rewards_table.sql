-- +goose Up
CREATE SCHEMA IF NOT EXISTS users;

CREATE TYPE users.reward_status AS ENUM ('active', 'redeemed');

CREATE TABLE IF NOT EXISTS users.reward_definitions (
    id   SERIAL PRIMARY KEY,
    code TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    earned_description TEXT NOT NULL,
    redeemed_description TEXT NOT NULL

);

INSERT INTO users.reward_definitions (
    code, name, description, earned_description, redeemed_description
)
VALUES
    (
        'DELIVERY_DISCOUNT_10',
        'Скидка 10% на доставку',
        'Скидка 10% на одну покупку с Авито Доставкой',
        'Получена скидка 10% на доставку',
        'Использована скидка 10% на доставку'
    ),
    (
        'FREE_LISTING_PROMOTION',
        'Бесплатное продвижение',
        'Одно бесплатное продвижение выбранного объявления',
        'Получено бесплатное продвижение объявления',
        'Использовано бесплатное продвижение объявления'
    ),
    (
        'AUTOTEKA_DISCOUNT_20',
        'Скидка 20% на Автотеку',
        'Скидка 20% на один отчёт об истории автомобиля',
        'Получена скидка 20% на отчёт Автотеки',
        'Использована скидка 20% на отчёт Автотеки'
    ),
    (
        'FREE_LISTING_HIGHLIGHT',
        'Выделение объявления',
        'Бесплатное визуальное выделение одного объявления на ограниченный срок',
        'Получено бесплатное выделение объявления',
        'Использовано бесплатное выделение объявления'
    ),
    (
        'LISTING_DISCOUNT_15',
        'Скидка 15% на размещение',
        'Скидка 15% на одно платное размещение объявления',
        'Получена скидка 15% на размещение объявления',
        'Использована скидка 15% на размещение объявления'
    );

CREATE TABLE IF NOT EXISTS users.user_rewards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users.users(id)
             ON DELETE CASCADE,
    reward_id INTEGER NOT NULL REFERENCES users.reward_definitions(id),
    promo_code TEXT NOT NULL UNIQUE,
    status users.reward_status NOT NULL DEFAULT 'active',
    earned_reason TEXT NOT NULL,
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
