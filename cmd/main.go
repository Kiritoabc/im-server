package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"im-system/internal/config"
	"im-system/internal/handler"
	"im-system/internal/middle"
	"im-system/internal/router"
	"im-system/internal/service"
)

func main() {
	// 初始化日志
	config.InitLogger()

	// 加载配置
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		config.Logger.Fatalf("加载配置失败: %v", err)
	}

	// 初始化MySQL
	err = config.InitDB(cfg)
	if err != nil {
		config.Logger.Fatalf("初始化MySQL失败: %v", err)
	}

	// 初始化Redis
	config.InitRedis(cfg)

	// 初始化服务和处理器
	userService := service.NewUserService()
	userHandler := handler.NewUserHandler(userService)

	friendService := service.NewFriendService()
	friendHandler := handler.NewFriendHandler(friendService)

	notificationService := service.NewNotificationService()
	notificationHandler := handler.NewNotificationHandler(notificationService)

	// 设置路由
	r := gin.Default()

	// 配置 CORS 中间件
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Authorization"},
	}))

	// 使用 JWT 中间件
	r.Use(middle.AuthMiddleware())

	router.RegisterRoutes(r, userHandler, friendHandler, notificationHandler)

	// 启动服务器
	config.Logger.Infof("HTTP服务器启动在端口%s\n", cfg.Server.HTTPPort)
	if err := r.Run(cfg.Server.HTTPPort); err != nil {
		config.Logger.Fatalf("启动服务器失败: %v", err)
	}
}
