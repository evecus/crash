package main

import (
	"crashpanel/api"
	"crashpanel/config"
	"crashpanel/database"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

//go:embed dist
var frontendDist embed.FS

func main() {
	// 加载应用配置
	cfg := config.Load()

	// 初始化数据库
	database.Init(cfg.DBPath)

	// 设置 Gin 模式
	if cfg.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// 注册 API 路由
	api.Register(r, cfg)

	// 内嵌前端静态文件
	distFS, err := fs.Sub(frontendDist, "dist")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load frontend dist: %v\n", err)
		os.Exit(1)
	}
	r.NoRoute(func(c *gin.Context) {
		// API 路由未匹配时返回前端 index.html（SPA 支持）
		if len(c.Request.URL.Path) > 4 && c.Request.URL.Path[:4] == "/api" {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		http.FileServer(http.FS(distFS)).ServeHTTP(c.Writer, c.Request)
	})

	addr := fmt.Sprintf(":%d", cfg.Port)
	fmt.Printf("CrashPanel running on http://0.0.0.0%s\n", addr)
	if err := r.Run(addr); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}
