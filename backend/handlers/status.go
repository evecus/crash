package handlers

import (
	"crashpanel/service/core"
	"crashpanel/service/system"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StatusHandler struct {
	manager *core.Manager
}

func NewStatusHandler(manager *core.Manager) *StatusHandler {
	return &StatusHandler{manager: manager}
}

func (h *StatusHandler) SystemInfo(c *gin.Context) {
	info := system.GetInfo()
	c.JSON(http.StatusOK, info)
}

func (h *StatusHandler) Network(c *gin.Context) {
	stats := h.manager.GetTraffic()
	c.JSON(http.StatusOK, stats)
}
