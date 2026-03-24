package handlers

import (
	"crashpanel/database"
	"crashpanel/models"
	"crashpanel/service/subscribe"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SubscriptionHandler struct{}

func NewSubscriptionHandler() *SubscriptionHandler {
	return &SubscriptionHandler{}
}

func (h *SubscriptionHandler) List(c *gin.Context) {
	var subs []models.Subscription
	database.DB.Order("created_at desc").Find(&subs)
	c.JSON(http.StatusOK, subs)
}

func (h *SubscriptionHandler) Create(c *gin.Context) {
	var sub models.Subscription
	if err := c.ShouldBindJSON(&sub); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sub.Status = "pending"
	if err := database.DB.Create(&sub).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sub)
}

func (h *SubscriptionHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var sub models.Subscription
	if err := database.DB.First(&sub, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	if err := c.ShouldBindJSON(&sub); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	database.DB.Save(&sub)
	c.JSON(http.StatusOK, sub)
}

func (h *SubscriptionHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	database.DB.Delete(&models.Subscription{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *SubscriptionHandler) Refresh(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var sub models.Subscription
	if err := database.DB.First(&sub, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	result, err := subscribe.Fetch(&sub)
	if err != nil {
		sub.Status = "error"
		sub.ErrorMsg = err.Error()
		database.DB.Save(&sub)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	sub.Status = "ok"
	sub.ErrorMsg = ""
	sub.NodeCount = result.NodeCount
	database.DB.Save(&sub)
	c.JSON(http.StatusOK, gin.H{
		"message":    "updated",
		"node_count": result.NodeCount,
	})
}
