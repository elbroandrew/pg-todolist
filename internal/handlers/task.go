package handlers

import (
	"errors"
	"net/http"
	"pg-todolist/internal/models"
	"pg-todolist/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
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
		if errors.Is(err, service.ErrTaskNotFound){

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

	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}
	
	userID := c.MustGet("userID").(uint)
	task.UserID = userID
	task.ID = uint(id)

	if err := h.taskService.Update(&task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
		return
	}

	c.JSON(http.StatusOK, task)

}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID задачи"})
		return
	}


	
	userID := c.MustGet("userID").(uint)


	if err := h.taskService.Delete(uint(id), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
		return
	}

	c.Status(http.StatusNoContent)

}