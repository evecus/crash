package handlers

import (
	"crashpanel/database"
	"crashpanel/models"
	"crashpanel/service/firewall"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FirewallHandler struct{}

func NewFirewallHandler() *FirewallHandler {
	return &FirewallHandler{}
}

func (h *FirewallHandler) GetRules(c *gin.Context) {
	var rules []models.FirewallRule
	database.DB.Order("created_at asc").Find(&rules)
	c.JSON(http.StatusOK, rules)
}

func (h *FirewallHandler) CreateRule(c *gin.Context) {
	var rule models.FirewallRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	database.DB.Create(&rule)
	c.JSON(http.StatusOK, rule)
}

func (h *FirewallHandler) DeleteRule(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	database.DB.Delete(&models.FirewallRule{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *FirewallHandler) Apply(c *gin.Context) {
	// 读取当前设置和规则
	var settings models.Settings
	if err := database.DB.First(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var rules []models.FirewallRule
	database.DB.Find(&rules)

	if err := firewall.Apply(&settings, rules); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "firewall rules applied"})
}

func (h *FirewallHandler) Flush(c *gin.Context) {
	var settings models.Settings
	database.DB.First(&settings)
	if err := firewall.Flush(settings.FirewallMod); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "firewall rules cleared"})
}
