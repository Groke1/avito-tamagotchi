## Тестирование

### GET /api/v1/tasks
```
curl.exe -X GET http://localhost:8081/api/v1/tasks -H "Authorization: Bearer $token" -H "X-User-ID: f8db7bab-1456-4061-a123-1bb3fad1abc3"

{"date":"2026-08-05","items":[{"id":"9b3541db-9942-4625-ba51-b6afcb23754a","title":"Лояльный продавец","description":"Получите новый отзыв с оценкой 5 звезд от верифицированного покупателя.","reward_coins":300,"reward_xp":500,"status":"active","completed_at":null},{"id":"4c202322-726f-4ec6-9896-62bc8b49d5ee","title":"Мастер описания","description":"Добавьте в активное объявление в категории \"Электроника\" видеообзор товара или подробные характеристики.","reward_coins":100,"reward_xp":200,"status":"active","completed_at":null},{"id":"e9b11c75-392c-473d-815a-52ef389d3110","title":"Первая продажа месяца","description":"Успешно завершите сделку по любому объявлению с подключенной Авито Доставкой.","reward_coins":500,"reward_xp":1000,"status":"active","completed_at":null}]}

```
### GET /api/v1/tasks/{task_id}
```
curl.exe -X GET http://localhost:8081/api/v1/tasks/9b3541db-9942-4625-ba51-b6afcb23754a -H "Authorization: Bearer $token" -H "X-User-ID: f8db7bab-1456-4061-a123-1bb3fad1abc3"

{"id":"9b3541db-9942-4625-ba51-b6afcb23754a","title":"Лояльный продавец","description":"Получите новый отзыв с оценкой 5 звезд от верифицированного покупателя.","reward_coins":300,"reward_xp":500,"status":"active","completed_at":null}

```
### POST /api/v1/tasks/{task_id}/complete
```
curl.exe -X POST http://localhost:8081/api/v1/tasks/9b3541db-9942-4625-ba51-b6afcb23754a/complete -H "Authorization: Bearer $token" -H "X-User-ID: f8db7bab-1456-4061-a123-1bb3fad1abc3"

{"task":{"id":"9b3541db-9942-4625-ba51-b6afcb23754a","title":"Лояльный продавец","description":"Получите новый отзыв с оценкой 5 звезд от верифицированного покупателя.","reward_coins":300,"reward_xp":500,"status":"completed","completed_at":"2026-08-05T16:26:34.386331502Z"},"awarded":{"coins":300,"xp":500}}
```
### плохие сценарии

```
curl.exe -X GET http://localhost:8081/api/v1/tasks -H "X-User-ID: f8db7bab-1456-4061-a123-1bb3fad1abc3"
{"code":"UNAUTHORIZED","message":"Требуется повторная авторизация"}

curl.exe -X GET http://localhost:8081/api/v1/tasks `
 -H "Authorization: Bearer $token"
{"code":"UNAUTHORIZED","message":"Идентификатор пользователя не найден"}

curl.exe -X GET http://localhost:8081/api/v1/tasks/00000000-0000-0000-0000-000000000000 `
 -H "Authorization: Bearer $token" `
 -H "X-User-ID: f8db7bab-1456-4061-a123-1bb3fad1abc3"
{"code":"TASK_NOT_FOUND","message":"Задание не найдено"}

curl.exe -X POST http://localhost:8081/api/v1/tasks/00000000-0000-0000-0000-000000000000/complete `
   -H "Authorization: Bearer $token" `
   -H "X-User-ID: f8db7bab-1456-4061-a123-1bb3fad1abc3"
{"code":"TASK_NOT_FOUND","message":"Задание не найдено"}

curl.exe -X POST http://localhost:8081/api/v1/tasks/9b3541db-9942-4625-ba51-b6afcb23754a/complete `
  -H "Authorization: Bearer $token" `
  -H "X-User-ID: f8db7bab-1456-4061-a123-1bb3fad1abc3"
{"code":"TASK_ALREADY_COMPLETED","message":"Награда за это задание уже получена"}

```