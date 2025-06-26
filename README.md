Используется паттерн Репозиторий - структурный паттерн для доступа к данным.



run: 

`docker-compose up -d`

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