package router

import (
	"github.com/gin-gonic/gin"
	"im-system/internal/handler"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(r *gin.Engine, userHandler *handler.UserHandler, friendHandler *handler.FriendHandler) {
	imGroup := r.Group("/im-server")
	// 注册登录模块
	{
		imGroup.POST("/register", userHandler.Register)
		imGroup.POST("/login", userHandler.Login)
		// user 模块
		imGroup.GET("/user/userInfo", userHandler.GetUserInfo)
		imGroup.POST("/user/:user_id/friend/:friend_id", userHandler.AddFriend)

		// 使用 friends 前缀
		friendsGroup := imGroup.Group("/friends")
		friendsGroup.POST("/user/:user_id/group", friendHandler.CreateFriendGroup) // 创建好友分组的路由
	}
}
