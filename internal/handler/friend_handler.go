package handler

import (
	"net/http"
	"strconv"

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

// CreateFriendGroup 创建好友分组
func (h *FriendHandler) CreateFriendGroup(c *gin.Context) {
	userIDStr := c.Param("user_id")
	groupName := c.Query("group_name") // 从查询参数获取分组名称

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		model.SendResponse(c, http.StatusBadRequest, model.Error("无效的用户ID"))
		return
	}

	if err := h.friendService.CreateFriendGroup(uint(userID), groupName); err != nil {
		config.Logger.Error(err)
		model.SendResponse(c, http.StatusInternalServerError, model.Error(err.Error()))
		return
	}

	model.SendResponse(c, http.StatusOK, model.Success("好友分组创建成功", nil))
}

// GetFriendGroups 获取用户的好友分组
func (h *FriendHandler) GetFriendGroups(c *gin.Context) {
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
