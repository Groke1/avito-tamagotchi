CREATE TYPE trip_status AS ENUM (
    'in_progress',
    'pending_delivery',
    'delivered'
);

CREATE TABLE IF NOT EXISTS trip_events (
    id SERIAL PRIMARY KEY,
    description TEXT NOT NULL,
    is_negative BOOLEAN NOT NULL
);

CREATE TABLE IF NOT EXISTS pet_trips (
    id BIGSERIAL PRIMARY KEY,
    pet_id BIGINT NOT NULL REFERENCES pets(id),
    user_id UUID NOT NULL,
    location VARCHAR(100) NOT NULL DEFAULT 'Zootopia',
    reward_xp INTEGER,
    reward_coins INTEGER,
    reward_code TEXT,
    story TEXT NOT NULL,
    status trip_status NOT NULL DEFAULT 'in_progress',
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ended_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_pet_trips_pet_id_started_at
    ON pet_trips (pet_id, started_at DESC);

CREATE INDEX idx_pet_trips_finished ON pet_trips (ended_at)
    WHERE status = 'in_progress';
;

INSERT INTO trip_events (description, is_negative) VALUES
    ('встретил доброго продавца', false),
    ('нашёл редкое объявление', false),
    ('купил вкусняшку', false),
    ('нашёл бесплатный диван', false),
    ('получил хорошую скидку', false),
    ('успел забрать товар', false),
    ('нашёл нужную запчасть', false),
    ('удачно поторговался', false),
    ('встретил старого друга', false),
    ('нашёл кладовку сокровищ', false),
    ('продал старый велосипед', false),
    ('нашёл винтажную игрушку', false),
    ('получил пять звёзд', false),
    ('забрал подарок даром', false),
    ('нашёл уютное кресло', false),
    ('купил почти даром', false),
    ('нашёл потерянный мячик', false),
    ('помог выбрать товар', false),
    ('встретил пушистого друга', false),
    ('нашёл новое хобби', false),

    ('продавец отменил встречу', true),
    ('товар уже продали', true),
    ('опоздал на встречу', true),
    ('заблудился по дороге', true),
    ('попал под дождь', true),
    ('промокли любимые лапки', true),
    ('забыл адрес встречи', true),
    ('продавец не ответил', true),
    ('не успел поторговаться', true),
    ('потерял любимый мячик', true),
    ('сел телефон', true),
    ('устал искать объявление', true),
    ('приехал не туда', true),
    ('застрял в пробке', true),
    ('товар оказался занят', true),
    ('убежал от дождя', true),
    ('испугался большой собаки', true),
    ('забыл взять монетки', true),
    ('пропустил нужный автобус', true),
    ('долго ждал продавца', true);