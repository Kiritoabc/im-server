package router

import (
	"im-system/internal/handler"
	"im-system/internal/module/kimi"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(r *gin.Engine, userHandler *handler.UserHandler,
	friendHandler *handler.FriendHandler,
	notificationHandler *handler.NotificationHandler,
	friendGroupHandler *handler.FriendGroupHandler,
	groupHandler *handler.GroupHandler,
	webSocketHandler *handler.WebSocketHandler,
	chatSummaryHandler *handler.ChatSummaryHandler) {

	// 配置静态文件服务
	r.Static("/static", "./static")

	imGroup := r.Group("/im-server")
	// 注册登录模块
	{
		imGroup.POST("/register", userHandler.Register)          // 注册接口
		imGroup.POST("/login", userHandler.Login)                // 登录接口
		imGroup.POST("/logout", userHandler.Logout)              // 退出登录接口
		imGroup.GET("/user/check_token", userHandler.CheckToken) // 检查token是否有效
		// user 模块
		imGroup.GET("/user/userInfo", userHandler.GetUserInfo)        // 获取用户信息
		imGroup.POST("/user/add_friend", userHandler.AddFriend)       // 添加好友的路由
		imGroup.POST("/user/update", userHandler.UpdateUserInfo)      // 更新用户信息
		imGroup.POST("/user/query", userHandler.QueryUserAndGroup)    // 查询用户和群聊信息
		imGroup.POST("/user/del_friend", userHandler.DeleteFriend)    // 删除好友
		imGroup.POST("/user/upload_avatar", userHandler.UploadAvatar) // 上传头像

		// 使用 friends 前缀
		friendsGroup := imGroup.Group("/friends")
		friendsGroup.GET("/user/groups", friendHandler.GetUserFriendAllFriends) // 获取好友分组的路由
		friendsGroup.GET("/user/friend_groups", friendHandler.GetFriendGroups)  // 获取用户的好友列表(用于好友模块)
		friendsGroup.GET("/usr/friends_chat", friendHandler.GetUserFriendsChat) // 获取用户的好友列表(用于私聊模块)
		friendsGroup.GET("/all", friendHandler.GetUserFriends)                  // 获取用户的好友列表(用于好友模块)
		friendsGroup.POST("/update", friendHandler.UpdateFriendGroup)           // 修改用户好友的分组

		// notifications 通知模块
		imGroup.GET("/notifications", notificationHandler.GetNotifications)                              // 获取通知的路由
		imGroup.POST("/notifications/:notification_id", notificationHandler.HandleFriendRequest)         // 处理好友请求
		imGroup.GET("/notifications/get/sent_notifications", notificationHandler.GetSentNotifications)   // 获取已发送的好友请求
		imGroup.GET("/notifications/friend_requests", notificationHandler.GetFriendRequestNotifications) // 获取好友请求通知

		// friend_groups 好友分组模块
		imGroup.POST("/friend_groups", friendGroupHandler.CreateFriendGroup)  // 创建好友分组
		imGroup.GET("/friend_groups", friendGroupHandler.GetUserFriendGroups) // 获取用户的所有好友分组

		// 群组模块
		imGroup.POST("/groups", groupHandler.CreateGroup)                         // 创建群组
		imGroup.POST("/groups/query", groupHandler.QueryGroups)                   // 查询群组
		imGroup.GET("/groups/user", groupHandler.GetUserGroups)                   // 获取用户所在的群聊
		imGroup.GET("/groups/my_groups", groupHandler.GetMyAllGroups)             // 获取用户的所有群聊
		imGroup.GET("/groups/members", groupHandler.GetGroupMembers)              // 获取群聊的所有成员
		imGroup.POST("/groups/invite", groupHandler.InviteGroup)                  // 邀请好友加入群聊
		imGroup.POST("/groups/update_member_role", groupHandler.UpdateMemberRole) // 更新群成员角色
		imGroup.POST("/groups/remove_member", groupHandler.RemoveMember)          // 移除群成员
		imGroup.POST("/groups/quit", groupHandler.QuitGroup)                      // 退出群聊
		imGroup.POST("/groups/update", groupHandler.UpdateGroup)                  // 更新群聊信息

		// WebSocket 路由，使用鉴权中间件
		imGroup.GET("/ws", kimi.HandleWebSocket)

		// 聊天模块
		imGroup.GET("/private/chat", webSocketHandler.SendMessage) // 发送消息

		// 添加新路由
		imGroup.POST("/chat/summary", chatSummaryHandler.HandleChatSummary)
	}
}
