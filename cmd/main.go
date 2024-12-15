package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"im-system/internal/config"
	"im-system/internal/handler"
	"im-system/internal/model"
	"im-system/internal/service"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化MySQL
	err = model.InitDB(
		cfg.MySQL.Host,
		cfg.MySQL.Port,
		cfg.MySQL.Username,
		cfg.MySQL.Password,
		cfg.MySQL.Database,
	)
	if err != nil {
		log.Fatalf("初始化MySQL失败: %v", err)
	}

	// 初始化MongoDB
	err = model.InitMongoDB(cfg.MongoDB.URI, cfg.MongoDB.Database)
	if err != nil {
		log.Fatalf("初始化MongoDB失败: %v", err)
	}

	// 初始化服务和处理器
	userService := service.NewUserService()
	userHandler := handler.NewUserHandler(userService)

	// 设置路由
	r := gin.Default()
	
	imGroup := r.Group("/im-server")
	{
		imGroup.POST("/register", userHandler.Register)
		imGroup.POST("/login", userHandler.Login)
		imGroup.POST("/add-friend", userHandler.AddFriend)
		imGroup.POST("/send-message", userHandler.SendMessage)
	}

	// 启动服务器
	log.Printf("HTTP服务器启动在端口%s\n", cfg.Server.HTTPPort)
	if err := r.Run(cfg.Server.HTTPPort); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
} 