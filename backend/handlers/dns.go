package handlers

import (
	"crashpanel/database"
	"crashpanel/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DNSHandler struct{}

func NewDNSHandler() *DNSHandler {
	return &DNSHandler{}
}

func (h *DNSHandler) Get(c *gin.Context) {
	var dns models.DNSConfig
	if err := database.DB.First(&dns).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dns)
}

func (h *DNSHandler) Update(c *gin.Context) {
	var dns models.DNSConfig
	if err := database.DB.First(&dns).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&dns); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	database.DB.Save(&dns)
	c.JSON(http.StatusOK, dns)
}
