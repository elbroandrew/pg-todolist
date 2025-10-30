package tests

import (
	"context"
	"fmt"
	"log"
	"os"
	"pg-todolist/internal/models"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"database/sql"

	gormysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	testDB *gorm.DB
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "mysql:latest",
		ExposedPorts: []string{"3306/tcp"},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "password",
			"MYSQL_DATABASE":      "test_db",
			"MYSQL_USER":          "test_user",
			"MYSQL_PASSWORD":      "test_password",
		},

		WaitingFor: wait.ForListeningPort("3306/tcp").WithStartupTimeout(180 * time.Second),
	}

	mysqlContainer, err := testcontainers.GenericContainer(ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})

	if err != nil {
		log.Fatalf("couldn't start testcontainer: %s", err)
	}

	defer func() {
		if err := mysqlContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()

	host, err := mysqlContainer.Host(ctx)
	if err != nil {
		log.Fatalf("failed to get container host: %v", err)
	}
	port, err := mysqlContainer.MappedPort(ctx, "3306")
	if err != nil {
		log.Fatalf("failed to get container port: %v", err)
	}

	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=30s", 
		"test_user", "test_password", host, port.Port(), "test_db")

	// Ждем пока база данных будет готова принимать подключения
	var db *sql.DB
	maxRetries := 10 // Увеличиваем количество попыток
	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("mysql", connStr)
		if err != nil {
			log.Printf("sql.Open failed (attempt %d): %v", i+1, err)
			time.Sleep(5 * time.Second)
			continue
		}

		err = db.PingContext(ctx)
		if err != nil {
			log.Printf("PingContext failed (attempt %d): %v", i+1, err)
			db.Close()
			time.Sleep(5 * time.Second)
			continue
		}
		log.Println("SUCCESS: database/sql connection")
		break
	}

	if err != nil {
		log.Fatalf("Failed to connect to database after %d attempts: %v", maxRetries, err)
	}
	defer db.Close()


	// Подключаемся через GORM
	log.Println("Attempting connection with GORM")
	gormDB, err := gorm.Open(gormysql.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatalf("GORM connection failed: %v", err)
	}
	log.Println("SUCCESS: GORM connection successful!")

	testDB = gormDB

	log.Println("Running migration on test database...")
	err = testDB.AutoMigrate(&models.User{}, &models.Task{})
	if err != nil {
		log.Fatalf("Failed to run auto-migrations on test DB: %v", err)
	}

	log.Println("Starting tests...")
	exitCode := m.Run()
	os.Exit(exitCode)
}

func clearTables(db *gorm.DB) {
	db.Exec("DELETE FROM tasks")
	db.Exec("DELETE FROM users")
	db.Exec("ALTER TABLE tasks AUTO_INCREMENT = 1")
	db.Exec("ALTER TABLE users AUTO_INCREMENT = 1")
}