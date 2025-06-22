run: 

`docker-compose up -d`

`docker-compose exec mysql sh -c 'mysql -u root -p todo_db < /app/migrations/001_init.sql'`

Get into the `todo_db` container:

`docker exec -ti mysql_todo mysql -u root -p todo_db`

`SHOW TABLES;`

For Go:

`go mod init pg-todolist`

`go mod tidy`