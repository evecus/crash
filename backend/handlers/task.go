package handlers

import (
	"crashpanel/database"
	"crashpanel/models"
	"crashpanel/service/task"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	scheduler *task.Scheduler
}

func NewTaskHandler(scheduler *task.Scheduler) *TaskHandler {
	return &TaskHandler{scheduler: scheduler}
}

func (h *TaskHandler) List(c *gin.Context) {
	var tasks []models.Task
	database.DB.Order("created_at asc").Find(&tasks)
	c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) Create(c *gin.Context) {
	var t models.Task
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := database.DB.Create(&t).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if t.Enabled {
		h.scheduler.Add(&t)
	}
	c.JSON(http.StatusOK, t)
}

func (h *TaskHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var t models.Task
	if err := database.DB.First(&t, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	database.DB.Save(&t)
	// 重新注册定时器
	h.scheduler.Remove(t.ID)
	if t.Enabled {
		h.scheduler.Add(&t)
	}
	c.JSON(http.StatusOK, t)
}

func (h *TaskHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	h.scheduler.Remove(uint(id))
	database.DB.Delete(&models.Task{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *TaskHandler) Run(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var t models.Task
	if err := database.DB.First(&t, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	log, code := h.scheduler.RunNow(&t)
	c.JSON(http.StatusOK, gin.H{
		"exit_code": code,
		"log":       log,
	})
}
