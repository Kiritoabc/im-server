package handler

import (
	"im-system/internal/model/db"
	"net/http"

	"github.com/gin-gonic/gin"
	"im-system/internal/config"
	"im-system/internal/model"
	"im-system/internal/model/dto"
	"im-system/internal/service"
)

type GroupHandler struct {
	groupService *service.GroupService
}

func NewGroupHandler(groupService *service.GroupService) *GroupHandler {
	return &GroupHandler{groupService: groupService}
}

// CreateGroup 创建群组
func (h *GroupHandler) CreateGroup(c *gin.Context) {
	var createGroupDTO dto.CreateGroupDTO
	if err := c.ShouldBindJSON(&createGroupDTO); err != nil {
		model.SendResponse(c, http.StatusBadRequest, model.Error("无效的请求"))
		return
	}

	// 从上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		model.SendResponse(c, http.StatusUnauthorized, model.Error("用户未登录"))
		return
	}
	// 从上下文中获取用户信息
	user, exists := c.Get("user_info")
	if !exists {
		model.SendResponse(c, http.StatusUnauthorized, model.Error("用户未登录"))
		return
	}
	// 将用户信息转换为User结构体
	userInfo, ok := user.(*db.User)
	if !ok {
		model.SendResponse(c, http.StatusUnauthorized, model.Error("用户信息错误"))
		return
	}

	group := db.Group{
		Name:        createGroupDTO.Name,
		OwnerID:     userID.(uint), // 设置群主的用户ID
		GroupAvatar: createGroupDTO.GroupAvatar,
	}

	if err := h.groupService.CreateGroup(group, *userInfo); err != nil {
		config.Logger.Error(err)
		model.SendResponse(c, http.StatusInternalServerError, model.Error(err.Error()))
		return
	}

	model.SendResponse(c, http.StatusOK, model.Success("群组创建成功", nil))
}

// QueryGroups 查询群组信息
func (h *GroupHandler) QueryGroups(c *gin.Context) {
	var queryGroupDTO dto.QueryGroupDTO
	if err := c.ShouldBindJSON(&queryGroupDTO); err != nil {
		model.SendResponse(c, http.StatusBadRequest, model.Error("无效的请求"))
		return
	}
	// todo: 参数校验
	groups, err := h.groupService.QueryGroups(queryGroupDTO.GroupID, queryGroupDTO.GroupName)
	if err != nil {
		config.Logger.Error(err)
		model.SendResponse(c, http.StatusInternalServerError, model.Error(err.Error()))
		return
	}

	model.SendResponse(c, http.StatusOK, model.Success("查询群组成功", groups))
}

// GetUserGroups 获取用户所在的群聊
func (h *GroupHandler) GetUserGroups(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		model.SendResponse(c, http.StatusUnauthorized, model.Error("用户未登录"))
		return
	}

	groups, err := h.groupService.GetUserGroups(userID.(uint))
	if err != nil {
		config.Logger.Error(err)
		model.SendResponse(c, http.StatusInternalServerError, model.Error(err.Error()))
		return
	}

	model.SendResponse(c, http.StatusOK, model.Success("获取群聊信息成功", groups))
}

// GetMyAllGroups 获取用户的所有群聊
func (h *GroupHandler) GetMyAllGroups(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		model.SendResponse(c, http.StatusUnauthorized, model.Error("用户未登录"))
		return
	}
	groups, err := h.groupService.GetMyAllGroups(userID.(uint))

	if err != nil {
		config.Logger.Error(err)
		model.SendResponse(c, http.StatusInternalServerError, model.Error("获取群聊失败"))
		return
	}
	model.SendResponse(c, http.StatusOK, model.Success("获取群聊信息成功", groups))
}
