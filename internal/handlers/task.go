package handlers

import (
	"errors"
	"net/http"
	"pg-todolist/internal/app_errors"
	"pg-todolist/internal/models"
	"pg-todolist/internal/service"
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
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}
	//get user ID from JWT
	userID := c.MustGet("userID").(uint)
	task.UserID = userID

	if err := h.taskService.Create(&task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
		return
	}

	c.JSON(http.StatusCreated, task)

}

func (h *TaskHandler) GetTasks(c *gin.Context) {
	//get user ID from JWT
	userID := c.MustGet("userID").(uint)
	tasks, err := h.taskService.GetTaskByUserID(userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
		return
	}

	c.JSON(http.StatusOK, tasks)

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

	c.JSON(http.StatusOK, task)

}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID задачи"})
		return
	}

	//беру только поле completed
	var request struct {
		Completed *bool `json:"completed"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}

	if request.Completed == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Поле completed обязательно!"})
		return
	}

	userID := c.MustGet("userID").(uint)

	// обновляю только нужные поля
	task := models.Task{
		Model: gorm.Model{ID: uint(id)},
		UserID:    userID,
		Completed: *request.Completed,
	}

	//Обновляю задачу
	rowsAffected, err := h.taskService.Update(&task)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Задача не найдена."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"task_id": id,
		"completed": *request.Completed,
		"rows_updated": rowsAffected,
	})

}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID задачи"})
		return
	}

	//помечает в БД как удалено, но физически удалит только так:
	//result := r.db.Unscoped().Delete(...)

	userID := c.MustGet("userID").(uint)

	// Проверяем существование задачи перед удалением
    if _, err := h.taskService.GetByID(uint(taskID), userID); err != nil {
        if errors.Is(err, app_errors.ErrTaskNotFound) {
            c.JSON(http.StatusNotFound, gin.H{
                "error": "Задача не найдена",
            })
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Ошибка при проверке задачи",
        })
        return
    }

	if _, err := h.taskService.Delete(uint(taskID), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
		return
	}

	c.Status(http.StatusNoContent)

}
