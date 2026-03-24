package handlers

import (
	"crashpanel/service/core"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CoreHandler struct {
	manager *core.Manager
}

func NewCoreHandler(manager *core.Manager) *CoreHandler {
	return &CoreHandler{manager: manager}
}

func (h *CoreHandler) Status(c *gin.Context) {
	c.JSON(http.StatusOK, h.manager.Status())
}

func (h *CoreHandler) Start(c *gin.Context) {
	if err := h.manager.Start(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "core started"})
}

func (h *CoreHandler) Stop(c *gin.Context) {
	if err := h.manager.Stop(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "core stopped"})
}

func (h *CoreHandler) Restart(c *gin.Context) {
	if err := h.manager.Restart(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "core restarted"})
}

func (h *CoreHandler) GetLog(c *gin.Context) {
	lines := 200
	if l := c.Query("lines"); l != "" {
		if n, err := strconv.Atoi(l); err == nil && n > 0 {
			lines = n
		}
	}
	c.JSON(http.StatusOK, gin.H{"lines": h.manager.GetLog(lines)})
}
