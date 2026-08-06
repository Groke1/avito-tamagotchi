-- +goose Up
-- +goose StatementBegin
INSERT INTO tasks (id, title, description, reward_coins, reward_xp) 
VALUES 
  (
    'e9b11c75-392c-473d-815a-52ef389d3110', 
    'Первая продажа месяца', 
    'Успешно завершите сделку по любому объявлению с подключенной Авито Доставкой.', 
    500, 1000
  ),
  (
    'a3e0b214-411a-463e-908c-94df1ba4e321', 
    'Быстрый ответ покупателю', 
    'Ответьте на 5 первых сообщений от разных покупателей в течение 15 минут с момента получения.', 
    150, 300
  ),
  (
    gen_random_uuid(), 
    'Мастер описания', 
    'Добавьте в активное объявление в категории "Электроника" видеообзор товара или подробные характеристики.', 
    100, 200
  ),
  (
    gen_random_uuid(), 
    'Лояльный продавец', 
    'Получите новый отзыв с оценкой 5 звезд от верифицированного покупателя.', 
    300, 500
  )
ON CONFLICT (id) DO NOTHING;

INSERT INTO user_tasks (user_id, task_id, status) 
VALUES 
  (
    '8c9c6123-289e-4e42-a89d-7db1387d81a9',
    'a3e0b214-411a-463e-908c-94df1ba4e321',
    'active'
  )
ON CONFLICT (user_id, task_id) DO NOTHING;

INSERT INTO user_tasks (user_id, task_id, status, completed_at) 
VALUES 
  (
    '8c9c6123-289e-4e42-a89d-7db1387d81a9', 
    'e9b11c75-392c-473d-815a-52ef389d3110',
    'completed', 
    NOW()
  )
ON CONFLICT (user_id, task_id) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM user_tasks 
WHERE user_id = '8c9c6123-289e-4e42-a89d-7db1387d81a9';

DELETE FROM tasks 
WHERE id IN (
  'e9b11c75-392c-473d-815a-52ef389d3110',
  'a3e0b214-411a-463e-908c-94df1ba4e321'
);
-- +goose StatementEnd