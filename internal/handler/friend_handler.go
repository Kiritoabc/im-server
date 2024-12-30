package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"

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
