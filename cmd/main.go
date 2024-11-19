package main

import (
	"audit-system/internal/api/router"
	"audit-system/pkg/config"
	"audit-system/pkg/database"
	"audit-system/pkg/logger"
	"fmt"

	"github.com/kataras/iris/v12"
)

func main() {
	// 初始化配置
	if err := config.Init(); err != nil {
		panic(fmt.Sprintf("配置初始化失败: %v", err))
	}

	// 初始化日志
	if err := logger.Init(); err != nil {
		panic(fmt.Sprintf("日志初始化失败: %v", err))
	}

	// 初始化数据库
	if err := database.InitDB(); err != nil {
		panic(fmt.Sprintf("数据库初始化失败: %v", err))
	}

	app := iris.New()

	// 注册路由
	router.RegisterRoutes(app)

	// 获取配置
	cfg := config.GetConfig()
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	// 启动服务器
	if err := app.Listen(addr); err != nil {
		panic(fmt.Sprintf("服务器启动失败: %v", err))
	}
}
