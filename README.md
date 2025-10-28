Проект состоит из двух микросервисов, один из которых отвечает за авторизацию пользователей с `JWT токеном`, реализует паттерн `API GATEWAY` через `Go Gin Reverse Proxy`,
другой микросервис `TaskService` отвечает за управление задачами пользователя. Оба сервиса соединяются с БД `MySQL`. Так же используется `Redis` для кеширования токенов - добавление их в `blacklist`.
Данный проект реализует паттерн `Репозиторий` - структурный паттерн для доступа к данным через сервисы и репозитории.
Имеется gui `https://github.com/elbroandrew/simple-todo-list`, написанный на чистом `js` без использования фреймоворка.


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