package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"im-system/internal/config"
	"im-system/internal/model"
	"im-system/internal/service"
)

type FriendHandler struct {
	friendService *service.FriendService
}

func NewFriendHandler(friendService *service.FriendService) *FriendHandler {
	return &FriendHandler{friendService: friendService}
}

// GetUserFriendAllFriends 获取用户的好友
func (h *FriendHandler) GetUserFriendAllFriends(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		model.SendResponse(c, http.StatusUnauthorized, model.Error("用户未登录"))
		return
	}

	groupVOs, err := h.friendService.GetFriendGroupsWithMembers(userID.(uint))
	if err != nil {
		config.Logger.Error(err)
		model.SendResponse(c, http.StatusInternalServerError, model.Error(err.Error()))
		return
	}

	model.SendResponse(c, http.StatusOK, model.Success("获取好友分组成功", groupVOs))
}

// GetUserFriendsChat 获取用户的好友聊天
func (h *FriendHandler) GetUserFriendsChat(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		model.SendResponse(c, http.StatusUnauthorized, model.Error("用户未登录"))
		return
	}
	friends, err := h.friendService.GetUserFriendsChat(userID.(uint))
	if err != nil {
		config.Logger.Error(err)
		model.SendResponse(c, http.StatusInternalServerError, model.Error(err.Error()))
		return
	}
	model.SendResponse(c, http.StatusOK, model.Success("获取好友列表成功", friends))
}

// GetUserFriends 获取用户的好友
func (h *FriendHandler) GetUserFriends(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		model.SendResponse(c, http.StatusUnauthorized, model.Error("用户未登录"))
		return
	}

	// 调用service层获取好友列表
	friends, err := h.friendService.GetUserFriends(userID.(uint))
	if err != nil {
		config.Logger.Error(err)
		model.SendResponse(c, http.StatusInternalServerError, model.Error(err.Error()))
		return
	}

	model.SendResponse(c, http.StatusOK, model.Success("获取好友列表成功", friends))
}

// GetFriendGroups 获取用户的所有好友分组
func (h *FriendHandler) GetFriendGroups(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		model.SendResponse(c, http.StatusUnauthorized, model.Error("用户未登录"))
		return
	}

	// 调用service层获取好友分组
	groups, err := h.friendService.GetUserFriendsGroups(userID.(uint))
	if err != nil {
		config.Logger.Error(err)
		model.SendResponse(c, http.StatusInternalServerError, model.Error(err.Error()))
		return
	}

	model.SendResponse(c, http.StatusOK, model.Success("获取好友分组成功", groups))
}

// SearchFriendGroups 搜索好友分组
func (h *FriendHandler) SearchFriendGroups(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		model.SendResponse(c, http.StatusUnauthorized, model.Error("用户未登录"))
		return
	}

	// 获取搜索关键词
	keyword := c.Query("keyword")
	if keyword == "" {
		model.SendResponse(c, http.StatusBadRequest, model.Error("请输入搜索关键词"))
		return
	}

	// 调用service层搜索好友分组
	groups, err := h.friendService.SearchFriendGroups(userID.(uint), keyword)
	if err != nil {
		config.Logger.Error(err)
		model.SendResponse(c, http.StatusInternalServerError, model.Error(err.Error()))
		return
	}

	model.SendResponse(c, http.StatusOK, model.Success("搜索好友分组成功", groups))
}
