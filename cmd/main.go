package main

import (
	"pg-todolist/pkg/database"
)

func main() {

	db:=database.InitMySQL()

	defer func(){
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()
	
}
