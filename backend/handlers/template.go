package handlers

import (
	"crashpanel/database"
	"crashpanel/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TemplateHandler struct{}

func NewTemplateHandler() *TemplateHandler { return &TemplateHandler{} }

func (h *TemplateHandler) List(c *gin.Context) {
	var templates []models.RuleTemplate
	database.DB.Order("id asc").Find(&templates)
	c.JSON(http.StatusOK, templates)
}

func (h *TemplateHandler) SetDefault(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	// 先清除所有默认
	database.DB.Model(&models.RuleTemplate{}).Where("1=1").Update("is_default", false)
	// 设置新默认
	if err := database.DB.Model(&models.RuleTemplate{}).Where("id = ?", id).Update("is_default", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已设为默认模板"})
}

func (h *TemplateHandler) Create(c *gin.Context) {
	var tmpl models.RuleTemplate
	if err := c.ShouldBindJSON(&tmpl); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	database.DB.Create(&tmpl)
	c.JSON(http.StatusOK, tmpl)
}

func (h *TemplateHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	database.DB.Delete(&models.RuleTemplate{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
