package handlers

import (
	"crashpanel/database"
	"crashpanel/models"
	"crashpanel/service/core"
	"crashpanel/service/subscribe"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SubscriptionHandler struct{}

func NewSubscriptionHandler() *SubscriptionHandler { return &SubscriptionHandler{} }

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
	if sub.LinkType == "" {
		sub.LinkType = models.LinkTypeURL
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

// Refresh 更新单个订阅，并重新生成配置文件
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

	if err := core.GenerateConfig(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message":    "订阅已更新，配置生成失败: " + err.Error(),
			"node_count": result.NodeCount,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":    "订阅已更新，配置已重新生成，重启核心后生效",
		"node_count": result.NodeCount,
	})
}

// RefreshAll 更新全部订阅，并重新生成配置文件
func (h *SubscriptionHandler) RefreshAll(c *gin.Context) {
	var subs []models.Subscription
	database.DB.Find(&subs)

	results := make([]gin.H, 0, len(subs))
	for _, sub := range subs {
		result, err := subscribe.Fetch(&sub)
		if err != nil {
			sub.Status = "error"
			sub.ErrorMsg = err.Error()
			database.DB.Save(&sub)
			results = append(results, gin.H{"id": sub.ID, "name": sub.Name, "error": err.Error()})
		} else {
			sub.Status = "ok"
			sub.ErrorMsg = ""
			sub.NodeCount = result.NodeCount
			database.DB.Save(&sub)
			results = append(results, gin.H{"id": sub.ID, "name": sub.Name, "node_count": result.NodeCount})
		}
	}

	// 所有订阅更新完后重新生成一次配置
	genErr := ""
	if err := core.GenerateConfig(); err != nil {
		genErr = err.Error()
	}

	c.JSON(http.StatusOK, gin.H{
		"results":      results,
		"config_error": genErr,
	})
}

// UploadFile 上传本地提供者文件
func (h *SubscriptionHandler) UploadFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file uploaded"})
		return
	}
	defer file.Close()

	dir := "/etc/crashpanel/core/providers"
	if err := os.MkdirAll(dir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	savePath := filepath.Join(dir, filepath.Base(header.Filename))
	out, err := os.Create(savePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "上传成功",
		"file_path": savePath,
	})
}

// GenerateConfig 手动触发生成配置文件
func (h *SubscriptionHandler) GenerateConfig(c *gin.Context) {
	if err := core.GenerateConfig(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "配置文件已重新生成，重启核心后生效"})
}
