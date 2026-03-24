package api

import (
	"crashpanel/config"
	"crashpanel/handlers"
	"crashpanel/middleware"
	"crashpanel/service/core"
	"crashpanel/service/task"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine, cfg *config.Config) {
	// CORS（开发时跨域）
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))

	// 初始化单例服务
	coreManager := core.NewManager()
	taskScheduler := task.NewScheduler()
	taskScheduler.Start()

	// 初始化 Handler
	authH := handlers.NewAuthHandler(cfg)
	settingsH := handlers.NewSettingsHandler()
	coreH := handlers.NewCoreHandler(coreManager)
	subH := handlers.NewSubscriptionHandler()
	fwH := handlers.NewFirewallHandler()
	dnsH := handlers.NewDNSHandler()
	taskH := handlers.NewTaskHandler(taskScheduler)
	statusH := handlers.NewStatusHandler(coreManager)

	api := r.Group("/api")

	// ── 公开路由（无需认证）──
	api.POST("/auth/login", authH.Login)

	// ── 需要认证的路由 ──
	auth := api.Group("/")
	auth.Use(middleware.Auth(cfg.JWTSecret))

	auth.GET("/auth/status", authH.Status)

	// 设置
	auth.GET("/settings", settingsH.Get)
	auth.PUT("/settings", settingsH.Update)

	// 核心管理
	auth.GET("/core/status", coreH.Status)
	auth.POST("/core/start", coreH.Start)
	auth.POST("/core/stop", coreH.Stop)
	auth.POST("/core/restart", coreH.Restart)
	auth.GET("/core/log", coreH.GetLog)

	// 订阅
	auth.GET("/subscriptions", subH.List)
	auth.POST("/subscriptions", subH.Create)
	auth.PUT("/subscriptions/:id", subH.Update)
	auth.DELETE("/subscriptions/:id", subH.Delete)
	auth.POST("/subscriptions/:id/refresh", subH.Refresh)

	// 防火墙
	auth.GET("/firewall/rules", fwH.GetRules)
	auth.POST("/firewall/rules", fwH.CreateRule)
	auth.DELETE("/firewall/rules/:id", fwH.DeleteRule)
	auth.POST("/firewall/apply", fwH.Apply)
	auth.POST("/firewall/flush", fwH.Flush)

	// DNS
	auth.GET("/dns", dnsH.Get)
	auth.PUT("/dns", dnsH.Update)

	// 计划任务
	auth.GET("/tasks", taskH.List)
	auth.POST("/tasks", taskH.Create)
	auth.PUT("/tasks/:id", taskH.Update)
	auth.DELETE("/tasks/:id", taskH.Delete)
	auth.POST("/tasks/:id/run", taskH.Run)

	// 系统信息
	auth.GET("/system/info", statusH.SystemInfo)
	auth.GET("/system/network", statusH.Network)
}
