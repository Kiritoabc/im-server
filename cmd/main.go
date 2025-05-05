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

const (
	// 定义常量，用于指定配置文件的路径
	dockerConfigPath = "config/config-docker.yaml"
	//localConfigPath  = "config/config-local.yaml"
	configPath = "config/config.yaml"
)

func main() {
	// 初始化日志
	config.InitLogger()
	config.Logger.Info("开始初始化")
	// 加载配置
	cfg, err := config.LoadConfig(dockerConfigPath)
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

	friendGroupService := service.NewFriendShipService()
	friendGroupHandler := handler.NewFriendShipHandler(friendGroupService)

	groupService := service.NewGroupService()
	groupHandler := handler.NewGroupHandler(groupService)

	// 初始化消息服务
	messageService := service.NewMessageService()
	chatSummaryHandler := handler.NewChatSummaryHandler(messageService)

	// 初始化 WebSocket 处理器
	webSocketHandler := handler.NewWebSocketHandler(messageService)
	//go webSocketHandler.StartMessageHandler() // 启动消息处理

	// 设置路由
	r := gin.Default()

	// 配置 CORS 中间件
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Authorization", "token"},
	}))

	// 使用 JWT 中间件
	r.Use(middle.AuthMiddleware())

	// 注册路由
	router.RegisterRoutes(r, userHandler, friendHandler, notificationHandler, friendGroupHandler, groupHandler, webSocketHandler, chatSummaryHandler)

	// 启动服务器
	config.Logger.Infof("HTTP服务器启动在端口%s\n", cfg.Server.HTTPPort)
	if err := r.Run(cfg.Server.HTTPPort); err != nil {
		config.Logger.Fatalf("启动服务器失败: %v", err)
	}
}
