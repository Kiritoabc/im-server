package router

import (
	"github.com/gin-gonic/gin"
	"im-system/internal/handler"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(r *gin.Engine, userHandler *handler.UserHandler, friendHandler *handler.FriendHandler, notificationHandler *handler.NotificationHandler) {
	imGroup := r.Group("/im-server")
	// 注册登录模块
	{
		imGroup.POST("/register", userHandler.Register)
		imGroup.POST("/login", userHandler.Login)
		// user 模块
		imGroup.GET("/user/userInfo", userHandler.GetUserInfo)
		imGroup.POST("/user/add_friend", userHandler.AddFriend) // 添加好友的路由

		// 使用 friends 前缀
		friendsGroup := imGroup.Group("/friends")
		friendsGroup.POST("/user/:user_id/group", friendHandler.CreateFriendGroup) // 创建好友分组的路由
		friendsGroup.GET("/user/groups", friendHandler.GetFriendGroups)            // 获取好友分组的路由

		// 通知模块
		imGroup.GET("/notifications", notificationHandler.GetNotifications)                      // 获取通知的路由
		imGroup.POST("/notifications/:notification_id", notificationHandler.HandleFriendRequest) // 处理好友请求
	}
}
