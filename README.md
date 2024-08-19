# Medods test task

Тестовое задание можно посмотреть в  **TODO.MD**

Микросервис для генерации токенов(Access и Refresh) по уникальному идентификатору пользовaтеля(GUID) и их обновления.

Используемые технологии:

+ Go
+ Postgresql
+ Gin
+ Docker-compose
+ Swagger
+ SMTP

Созданный сервис имеет чистую архитектуру, что обеспечивает простое расширение его возможностей и удобное тестирование. Также в нем реализован Graceful Shutdown для правильного завершения работы сервиса.

# Environment

Создайте `.env` файл с нужными переменными  или переименуйте `.env.example`:
```.env
SERVER_PORT=7777
DB_USER=postgres
DB_PASSWORD=postgres
DB_HOST=127.0.0.1
DB_PORT=5432
DB_NAME=medods
ACCESS_KEY=medodskey
ACCESS_TOKEN_AGE=1h
REFRESH_TOKEN_AGE=720h
SMTP_USERNAME= ваша почта
SMTP_PASSWORD=asdg sdfd aeas hfgh  ПРИМЕР
SMTP_PORT=587
SMTP_SERVER=smtp.gmail.com
SMTP_RECEIVER= почта куда нужно отправить уведомление
```
Так же можно воспользоваться тестовым  `.env.test` переименовав в `.env`

# Usage

Если вы хотите запустить через compose,то:

* Для запуска контейнера:

```shell
make start
```
* Для остановки контейнера:

```shell
make stop
```

Если вручную:

* Запустить докер postgres:

```shell
make postgres
```

* Создать бд:

```shell
make createdb
```


* Затем запуск:

```shell
make run
```

По пути http://0.0.0.0:7777/swagger/index.html сваггер документация

# Routes

* POST /api/v1/auth/login
  * Возвращает пару access & refresh токенов при вводе guid в тело запроса при валидном guid
  # Тело запроса
  ```json
    {
        "guid":"388d6b48-1a7e-4df7-b1eb-79ca726fb814"
    }
  ```
  # Тело ответа
  ```json
    {
        "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDc0ODE1MzAsInNlc3Npb25faWQiOiIyZjU1YjQ2ZS1lZTEyLTQ1MDUtYTkxNC1lZjRlOWYyMjlkYWQiLCJndWlkIjoiMzg4ZDZiNDgtMWE3ZS00ZGY3LWIxZWItNzljYTcyNmZiODE0In0.ChSsvpfPXVbtc7B8N1kaRn7_8FtIHqN0Am6OI9ADgM6f9CeER1T73xwHmfCOppv6bd6AUowj7lMe32dp3AIPWg",
        "refresh_token": "DAPNlaA8R3mgktccJwsy1g=="
    }
  ```
* POST /api/v1/auth/refresh
  * Возвращает при вводе в тело запроса session_id пару обновленных access & refresh token
  ```json
    {
        "session_id":"2f55b46e-ee12-4505-a914-ef4e9f229dad"
    }
  ```
  # Тело ответа
  ```json
    {
        "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDc0ODE5MDIsInNlc3Npb25faWQiOiIyZjU1YjQ2ZS1lZTEyLTQ1MDUtYTkxNC1lZjRlOWYyMjlkYWQiLCJndWlkIjoiMzg4ZDZiNDgtMWE3ZS00ZGY3LWIxZWItNzljYTcyNmZiODE0In0.mP2Vbzgtqv_0hnFt5V050Xnyqg9kB_zSL2KWkmY7NsMLSxqGhEd7B1P0K_jq3IG64rmYHQR1rFygoAmkl8fQgQ",
        "refresh_token": "iSgH6v0USyOqIui1sLexZA=="
    }
  ```
Дополнительно добавил route getAllSessions для просмотра всех сессий.

* GET /api/v1/auth/sessions
    * Возвращает все сессии
     # Тело ответа
    ```json
    [
        {
            "id": "6922e38d-6727-430f-ad02-a580c3087464",
            "guid": "9b73b94b-ddbf-4008-a391-1d5459491ca5",
            "hashed_refresh_token": "JDJhJDEwJFFVTTRsRjFXU2JmTnRBcFlXRzVocC5MQ1B4UUdOejd1TnNkNThsanZsWXJQL3JEbXdsNzVl",
            "created_time": "2024-02-09T11:59:02.115Z",
            "updated_time": "2024-02-09T12:00:01.14Z"
        },
        {
            "id": "2f55b46e-ee12-4505-a914-ef4e9f229dad",
            "guid": "388d6b48-1a7e-4df7-b1eb-79ca726fb814",
            "hashed_refresh_token": "JDJhJDEwJGZSUHdIWEN4UUI3L3pGWHVaZTVTLi5LOFVmNlYxTWp5Z1Q2alk2QXlGVWozZzNkZEN1QnRD",
            "created_time": "2024-02-09T12:25:30.787Z",
            "updated_time": "2024-02-09T12:31:42.985Z"
        }
    ]
    ```


