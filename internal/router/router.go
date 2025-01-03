package router

import (
	"github.com/gin-gonic/gin"
	"im-system/internal/handler"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(r *gin.Engine, userHandler *handler.UserHandler,
	friendHandler *handler.FriendHandler,
	notificationHandler *handler.NotificationHandler,
	friendGroupHandler *handler.FriendGroupHandler) {
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
		friendsGroup.GET("/user/groups", friendHandler.GetUserFriendAllFriends) // 获取好友分组的路由

		// notifications 通知模块
		imGroup.GET("/notifications", notificationHandler.GetNotifications)                      // 获取通知的路由
		imGroup.POST("/notifications/:notification_id", notificationHandler.HandleFriendRequest) // 处理好友请求

		// friend_groups 好友分组模块
		imGroup.POST("/friend_groups", friendGroupHandler.CreateFriendGroup)  // 创建好友分组
		imGroup.GET("/friend_groups", friendGroupHandler.GetUserFriendGroups) // 获取用户的所有好友分组

		// 获取用户发出的所有通知请求
		imGroup.GET("/sent_notifications", notificationHandler.GetSentNotifications) // 获取发出的通知请求
	}
}
