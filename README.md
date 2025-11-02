Проект состоит из двух микросервисов, один из которых отвечает за авторизацию пользователей с `JWT токеном`, реализует паттерн `API GATEWAY` через `Go Gin Reverse Proxy`,
другой микросервис `TaskService` отвечает за управление задачами пользователя. Оба сервиса соединяются с БД `MySQL`. Так же используется `Redis` для кеширования токенов - добавление их в `blacklist`.
Данный проект реализует паттерн `Репозиторий` - структурный паттерн для доступа к данным через сервисы и репозитории.
Имеется gui `https://github.com/elbroandrew/simple-todo-list`, написанный на чистом `js` без использования фреймоворка.

`cmd/gateway` - основной шлюз, 

`cmd/task_service` - сервис задач, 

`cmd/migrator` - мигратор для БД.

Проект имеет юнит и интеграционные тесты, которые запускаются с помощью `testcontainers-go`.

Описание пайплайна:

`docker-compose -d up` - запускает `mysql` и `redis`. Ждёт, пока `mysql` пройдёт `healthcheck` и станет `healthy`.
Далее запускает контейнер `migrator`. Он подключается к БД, применяет авто-миграции и завершает работу. Docker compose видит, что мигратор завершился успешно и запускает `task_service`, т.к. его зависимости `mysql` и `migrator` удовлетворены. Ждёт пока `task_service` станет `healthy` и запускает `gateway`.

`Reverse Proxy` путь:

1. Клиент отправляет `POST` `/tasks` с `JWT токеном`.

2. `Gateway`: проверяет токен, извлекает `userID=1`.

3. `Reverse Proxy`: создает запрос к `http://task-service:8080/` с заголовком `X-User-ID: 1`.

4. `TaskService`: получает запрос, видит `X-User-ID: 1`, создает задачу для пользователя 1.

5. `Reverse Proxy`: получает ответ, копирует его клиенту.


run: 

`docker-compose up -d`

Не обязательно, т.к. автомиграции есть

`docker-compose exec mysql sh -c 'mysql -u root -p todo_db < /app/migrations/001_init.sql'`

Get into the `todo_db` container:

`docker exec -ti mysql_todo mysql -u root -p todo_db`

MySQL:

`SHOW TABLES;`

Или показать схемы таблиц:

`docker-compose exec mysql sh -c 'mysql -u root -p -e "USE todo_db; DESCRIBE users; DESCRIBE tasks;"'`

Redis показать все ключи:

`docker exec -ti redis_todo redis-cli`

`KEYS *`

`FLUSHDB`

For Go:

`go mod init pg-todolist`

`go mod tidy`

`go run cmd/main.go`

Запуск тестов из корня проекта:

`go test ./...`

Запуск интеграционных тестов:

`go test -v ./tests/...`

# TODO:

- Добавить еще один микросервис: NotificationService

- Добавить gRPC

- Circuit breaker (для NotificationService) - Graceful degradation - система деградирует постепенно, а не падает полностью

- Frontend + nginx (и оставить реверс прокси пока что)

- CI/CD с помощью GitHub Actions

- Добавить Мониторинг и Логирование
