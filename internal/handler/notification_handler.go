package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"im-system/internal/config"
	"im-system/internal/model"
	"im-system/internal/service"
)

type NotificationHandler struct {
	notificationService *service.NotificationService
}

func NewNotificationHandler(notificationService *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationService: notificationService}
}

// GetNotifications 获取用户的通知
func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		model.SendResponse(c, http.StatusUnauthorized, model.Error("用户未登录"))
		return
	}

	// 获取 Type 参数
	notificationType := c.Query("type")

	notifications, err := h.notificationService.GetNotifications(userID.(uint), notificationType)
	if err != nil {
		config.Logger.Error(err)
		model.SendResponse(c, http.StatusInternalServerError, model.Error(err.Error()))
		return
	}

	model.SendResponse(c, http.StatusOK, model.Success("获取通知成功", notifications))
}

// HandleFriendRequest 处理好友请求
func (h *NotificationHandler) HandleFriendRequest(c *gin.Context) {
	notificationIDStr := c.Param("notification_id")
	action := c.Query("action") // "accept" 或 "reject"
	groupIdStr := c.Query("group_id")

	notificationID, err := strconv.ParseUint(notificationIDStr, 10, 32)
	if err != nil {
		model.SendResponse(c, http.StatusBadRequest, model.Error("无效的通知ID"))
		return
	}
	groupId, err := strconv.ParseUint(groupIdStr, 10, 32)
	if err != nil {
		model.SendResponse(c, http.StatusBadRequest, model.Error("无效的通知群组ID"))
		return
	}
	if action == "accept" {
		if err := h.notificationService.AcceptFriendRequest(uint(notificationID), uint(groupId)); err != nil {
			config.Logger.Error(err)
			model.SendResponse(c, http.StatusInternalServerError, model.Error(err.Error()))
			return
		}
		model.SendResponse(c, http.StatusOK, model.Success("好友请求已接受", nil))
	} else if action == "reject" {
		if err := h.notificationService.RejectFriendRequest(uint(notificationID)); err != nil {
			config.Logger.Error(err)
			model.SendResponse(c, http.StatusInternalServerError, model.Error(err.Error()))
			return
		}
		model.SendResponse(c, http.StatusOK, model.Success("好友请求已拒绝", nil))
	} else {
		model.SendResponse(c, http.StatusBadRequest, model.Error("无效的操作"))
	}
}
