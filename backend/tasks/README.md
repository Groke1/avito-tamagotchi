## Тестирование

### GET /api/v1/tasks

**Как работает?** Берется рандомное кол-во задач (3-5) и логируется в таблицу user_task. 
При повторном запросе возвращаются уже сгенерированные задачи. Задачи обновляются каждый день

```
>> $body = @{
 username = "kot_master"
 email    = "kot@avito.ru"
    password = "supersecret123"
 } | ConvertTo-Json

>> $response = Invoke-RestMethod -Method POST -Uri "http://localhost:8080/api/v1/auth/register" -ContentType "application/json" -Body $body

>> $accessToken= $response.access_token

>> $r = Invoke-RestMethod     -Method GET `
-Uri "http://localhost:8081/api/v1/tasks" `
-Headers @{
Authorization = "Bearer $accessToken"}

>> $r | ConvertTo-Json -Depth 10
{
    "date":  "2026-08-07",
    "items":  [
                  {
                      "id":  "cd51eda4-c593-4037-a0c1-c5bf3e7da86f",
                      "title":  "Лояльный продавец",
                      "description":  "Получите новый отзыв с оценкой 5 звезд от верифицированного покупателя.",
                      "reward_coins":  300,
                      "reward_xp":  500,
                      "status":  "active",
                      "completed_at":  null
                  },
                  {
                      "id":  "a3e0b214-411a-463e-908c-94df1ba4e321",
                      "title":  "Быстрый ответ покупателю",
                      "description":  "Ответьте на 5 первых сообщений от разных покупателей в течение 15 минут с момента получения.",
                      "reward_coins":  150,
                      "reward_xp":  300,
                      "status":  "active",
                      "completed_at":  null
                  },
                  {
                      "id":  "c1aa015e-c7f4-4ecb-807d-c05489c1a460",
                      "title":  "Мастер описания",
                      "description":  "Добавьте в активное объявление в категории \"Электроника\" видеообзор товара или подробные характеристики.",
                      "reward_coins":  100,
                      "reward_xp":  200,
                      "status":  "active",
                      "completed_at":  null
                  },
                  {
                      "id":  "e9b11c75-392c-473d-815a-52ef389d3110",
                      "title":  "Первая продажа месяца",
                      "description":  "Успешно завершите сделку по любому объявлению с подключенной Авито Доставкой.",
                      "reward_coins":  500,
                      "reward_xp":  1000,
                      "status":  "active",
                      "completed_at":  null
                  }
              ]
}
```
### GET /api/v1/tasks/{task_id}

```
$taskId = "cd51eda4-c593-4037-a0c1-c5bf3e7da86f"

curl.exe -X GET "http://localhost:8081/api/v1/tasks/$taskId" -H "Authorization: Bearer $($accessToken)"


{
   "id":"cd51eda4-c593-4037-a0c1-c5bf3e7da86f",
   "title":"Лояльный продавец",
   "description":"Получите новый отзыв с оценкой 5 звезд от верифицированного покупателя.",
   "reward_coins":300,
   "reward_xp":500,
   "status":"completed",
   "completed_at":"2026-08-07T08:56:50.215148Z"
}
```
### POST /api/v1/tasks/{task_id}/complete
```
$headers = @{
   "Authorization" = "Bearer $accessToken"
}

$taskId = "cd51eda4-c593-4037-a0c1-c5bf3e7da86f"

Invoke-RestMethod -Method POST -Uri "http://localhost:8081/api/v1/tasks/$taskId/complete" -Headers $headers | ConvertTo-Json -Depth 10


{"task":{"id":"9b3541db-9942-4625-ba51-b6afcb23754a","title":"Лояльный продавец","description":"Получите новый отзыв с оценкой 5 звезд от верифицированного покупателя.","reward_coins":300,"reward_xp":500,"status":"completed","completed_at":"2026-08-05T16:26:34.386331502Z"},"awarded":{"coins":300,"xp":500}}
```
