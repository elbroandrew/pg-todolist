package handlers

import (
	"errors"
	"net/http"
	"pg-todolist/internal/app_errors"
	"pg-todolist/internal/models"
	"pg-todolist/internal/service"
	"pg-todolist/internal/dto"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TaskHandler struct {
	taskService *service.TaskService
}

func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req dto.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}
	//get user ID from JWT
	userID := c.MustGet("userID").(uint)
	task := &models.Task{
		Title: req.Title,
		UserID: userID,
	}

	if err := h.taskService.Create(task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания задачи."})
		return
	}

	c.JSON(http.StatusCreated, dto.TaskResponseFromModel(task))

}

func (h *TaskHandler) GetTasks(c *gin.Context) {
	//get user ID from JWT
	userID := c.MustGet("userID").(uint)
	tasks, err := h.taskService.GetTaskByUserID(userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tasks."})
		return
	}

	c.JSON(http.StatusOK, dto.TasksResponseFromModels(tasks))

}

func (h *TaskHandler) GetTaskByID(c *gin.Context) {
	// Get ID from a URL parameter
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID"})
		return
	}

	userID := c.MustGet("userID").(uint)

	task, err := h.taskService.GetByID(uint(id), userID)
	if err != nil {
		if errors.Is(err, app_errors.ErrTaskNotFound) {

			c.JSON(http.StatusNotFound, gin.H{"error": "Задача не найдена"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
		return
	}

	c.JSON(http.StatusOK, dto.TaskResponseFromModel(task))

}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID задачи"})
		return
	}

	var req dto.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":"Invalid request data."})
		return
	}

	userID := c.MustGet("userID").(uint)

	//Обновляю задачу
	updatedTask, err := h.taskService.Update(uint(id), userID, *req.Completed)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, app_errors.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Задача не найдена."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task."})
		return
	}

	c.JSON(http.StatusOK, dto.TaskResponseFromModel(updatedTask))

}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID задачи"})
		return
	}

	//помечает в БД как удалено, но физически удалит только так:
	//result := r.db.Unscoped().Delete(...)

	userID := c.MustGet("userID").(uint)

	if err := h.taskService.Delete(uint(id), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
		return
	}

	c.Status(http.StatusNoContent)

}
