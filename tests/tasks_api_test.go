package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"pg-todolist/internal/handlers"
	"pg-todolist/internal/models"
	"pg-todolist/internal/repository"
	"pg-todolist/internal/router"
	"pg-todolist/internal/service"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)


func setUpTestApp() *gin.Engine{

	mockTokenService := new(service.TokenServiceMock)

	/*
	когда AuthMiddleware вызовет ValidateAccessToken,
	мок должен вернуть userID 1 и отсутствие ошибки.
	mock.Anything - потому что не важен сам токен, его сгенерирую позже.
	*/
	mockTokenService.On("ValidateAccessToken", mock.Anything).Return(1, nil)

	userRepo := repository.NewUserRepository(testDB)
	taskRepo := repository.NewTaskRepository(testDB)

	jwtSecret := "secret123"
	os.Setenv("JWT_SECRET", jwtSecret)

	taskService := service.NewTaskService(taskRepo)
	authService := service.NewAuthService(userRepo)

	taskHandler := handlers.NewTaskHandler(taskService)
	authHandler := handlers.NewAuthHandler(authService, mockTokenService)

	taskServiceRouter := router.SetupTaskServiceRouter(taskHandler)
	proxyTargetServer := httptest.NewServer(taskServiceRouter)
	gatewayRouter := router.SetupGatewayRouter(authHandler, mockTokenService, nil, proxyTargetServer.URL)

	return gatewayRouter
}


func TestTasksAPI_CreateAndGetTasks(t *testing.T){
	clearTables(testDB)

	
	testRouter := setUpTestApp()

	//new user
	testUser := &models.User{
		Email: "test_user@test.com",
		Password: "pa$$w0rd1",
	}

	testDB.Create(testUser) // ID 1
	jwtSecret := "secret123"
	os.Setenv("JWT_SECRET", jwtSecret)

	tokenService := service.NewTokenService(jwtSecret)

	accessToken, _, _ := tokenService.GenerateTokenPair(testUser.ID)
	authHeader := "Bearer " + accessToken

	// new task
	taskPayload := `{"title": "Integration test 1"}`
	reqCreate, _ := http.NewRequest("POST", "/tasks", bytes.NewBufferString(taskPayload))
	reqCreate.Header.Set("Content-Type", "application/json")
	reqCreate.Header.Set("Authorization", authHeader)

	wCreate := httptest.NewRecorder()
	testRouter.ServeHTTP(wCreate, reqCreate)


	//assert
	assert.Equal(t, http.StatusCreated, wCreate.Code, "Expected status 201 Created")

	var createdTask models.Task
	err := json.Unmarshal(wCreate.Body.Bytes(), &createdTask)
	assert.NoError(t, err)
	assert.Equal(t, "Integration test 1", createdTask.Title)

	// get all tasks
	reqGet, _ := http.NewRequest("GET", "/tasks", nil)
	reqGet.Header.Set("Authorization", authHeader)

	wGet := httptest.NewRecorder()
	testRouter.ServeHTTP(wGet, reqGet)

	assert.Equal(t, http.StatusOK, wGet.Code, "Expected status 200 OK")

	var tasks []models.Task
	err = json.Unmarshal(wGet.Body.Bytes(), &tasks)
	assert.NoError(t, err)

	//check if lists contain 1 task
	assert.Len(t, tasks, 1, "Expected to get one task")
	assert.Equal(t, createdTask.ID, tasks[0].ID)
	assert.Equal(t, "Integration test 1", tasks[0].Title)


}