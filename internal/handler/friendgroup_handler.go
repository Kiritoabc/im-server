package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"im-system/internal/config"
	"im-system/internal/model"
	"im-system/internal/model/db"
	"im-system/internal/service"
)

type FriendGroupHandler struct {
	friendShipService *service.FriendGroupService
}

func NewFriendShipHandler(friendShipService *service.FriendGroupService) *FriendGroupHandler {
	return &FriendGroupHandler{friendShipService: friendShipService}
}

// GetUserFriendGroups 获取当前用户的所有好友分组
func (h *FriendGroupHandler) GetUserFriendGroups(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		model.SendResponse(c, http.StatusUnauthorized, model.Error("用户未登录"))
		return
	}

	groups, err := h.friendShipService.GetFriendGroupsWithMembers(userID.(uint))
	if err != nil {
		config.Logger.Error(err)
		model.SendResponse(c, http.StatusInternalServerError, model.Error(err.Error()))
		return
	}

	model.SendResponse(c, http.StatusOK, model.Success("获取好友分组成功", groups))
}

// CreateFriendGroup 创建好友分组
func (h *FriendGroupHandler) CreateFriendGroup(c *gin.Context) {
	var group db.FriendGroup
	if err := c.ShouldBindJSON(&group); err != nil {
		model.SendResponse(c, http.StatusBadRequest, model.Error("无效的请求"))
		return
	}

	// 从上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		model.SendResponse(c, http.StatusUnauthorized, model.Error("用户未登录"))
		return
	}

	group.UserID = userID.(uint) // 设置分组的用户ID

	if err := h.friendShipService.CreateFriendGroup(group.UserID, group.GroupName); err != nil {
		config.Logger.Error(err)
		model.SendResponse(c, http.StatusInternalServerError, model.Error(err.Error()))
		return
	}

	model.SendResponse(c, http.StatusOK, model.Success("好友分组创建成功", nil))
}
