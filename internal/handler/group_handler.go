package handler

import (
	"im-system/internal/model/db"
	"net/http"

	"im-system/internal/config"
	"im-system/internal/model"
	"im-system/internal/model/dto"
	"im-system/internal/service"

	"github.com/gin-gonic/gin"
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

// GetGroupMembers 获取群聊成员
func (h *GroupHandler) GetGroupMembers(c *gin.Context) {
	// 获取参数
	groupID := c.Query("group_id")
	if groupID == "" {
		model.SendResponse(c, http.StatusBadRequest, model.Error("无效的请求"))
		return
	}
	members, err := h.groupService.GetGroupMembers(groupID)
	if err != nil {
		config.Logger.Error(err)
		model.SendResponse(c, http.StatusInternalServerError, model.Error("获取群成员失败"))
	}

	model.SendResponse(c, http.StatusOK, model.Success("获取群成员成功", members))
}

// InviteGroup 邀请好友加入群聊
func (h *GroupHandler) InviteGroup(c *gin.Context) {
	var inviteGroupDTO dto.InviteGroupDTO
	if err := c.ShouldBindJSON(&inviteGroupDTO); err != nil {
		model.SendResponse(c, http.StatusBadRequest, model.Error("无效的请求"))
		return
	}

	// 从上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		model.SendResponse(c, http.StatusUnauthorized, model.Error("用户未登录"))
		return
	}

	// 调用service层处理邀请逻辑
	if err := h.groupService.InviteGroup(userID.(uint), inviteGroupDTO.GroupID, inviteGroupDTO.FriendIDs); err != nil {
		config.Logger.Error(err)
		model.SendResponse(c, http.StatusInternalServerError, model.Error("邀请失败"))
		return
	}

	model.SendResponse(c, http.StatusOK, model.Success("邀请发送成功", nil))
}

// UpdateMemberRole 更新群成员角色
func (h *GroupHandler) UpdateMemberRole(c *gin.Context) {
	var updateRoleDTO struct {
		GroupID  uint   `json:"group_id" binding:"required"`
		MemberID uint   `json:"member_id" binding:"required"`
		Role     string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&updateRoleDTO); err != nil {
		model.SendResponse(c, http.StatusBadRequest, model.Error("无效的请求"))
		return
	}

	// 从上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		model.SendResponse(c, http.StatusUnauthorized, model.Error("用户未登录"))
		return
	}

	// 检查当前用户是否有权限修改成员角色
	var currentMember db.GroupMember
	if err := h.groupService.GetDB().Where("group_id = ? AND user_id = ?", updateRoleDTO.GroupID, userID).First(&currentMember).Error; err != nil {
		model.SendResponse(c, http.StatusForbidden, model.Error("您不是该群组成员"))
		return
	}

	// 只有群主和管理员可以修改成员角色
	if currentMember.Role != "owner" && currentMember.Role != "admin" {
		model.SendResponse(c, http.StatusForbidden, model.Error("您没有权限修改成员角色"))
		return
	}

	// 管理员不能修改群主角色
	if currentMember.Role == "admin" {
		var targetMember db.GroupMember
		if err := h.groupService.GetDB().Where("group_id = ? AND user_id = ?", updateRoleDTO.GroupID, updateRoleDTO.MemberID).First(&targetMember).Error; err == nil {
			if targetMember.Role == "owner" {
				model.SendResponse(c, http.StatusForbidden, model.Error("管理员不能修改群主角色"))
				return
			}
		}
	}

	// 更新成员角色
	if err := h.groupService.UpdateMemberRole(updateRoleDTO.GroupID, updateRoleDTO.MemberID, updateRoleDTO.Role); err != nil {
		config.Logger.Error(err)
		model.SendResponse(c, http.StatusInternalServerError, model.Error(err.Error()))
		return
	}

	model.SendResponse(c, http.StatusOK, model.Success("更新成员角色成功", nil))
}

// RemoveMember 移除群成员
func (h *GroupHandler) RemoveMember(c *gin.Context) {
	var removeMemberDTO struct {
		GroupID  uint `json:"group_id" binding:"required"`
		MemberID uint `json:"member_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&removeMemberDTO); err != nil {
		model.SendResponse(c, http.StatusBadRequest, model.Error("无效的请求"))
		return
	}

	// 从上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		model.SendResponse(c, http.StatusUnauthorized, model.Error("用户未登录"))
		return
	}

	// 检查当前用户是否有权限移除成员
	var currentMember db.GroupMember
	if err := h.groupService.GetDB().Where("group_id = ? AND user_id = ?", removeMemberDTO.GroupID, userID).First(&currentMember).Error; err != nil {
		model.SendResponse(c, http.StatusForbidden, model.Error("您不是该群组成员"))
		return
	}

	// 只有群主和管理员可以移除成员
	if currentMember.Role != "owner" && currentMember.Role != "admin" {
		model.SendResponse(c, http.StatusForbidden, model.Error("您没有权限移除成员"))
		return
	}

	// 管理员不能移除其他管理员
	if currentMember.Role == "admin" {
		var targetMember db.GroupMember
		if err := h.groupService.GetDB().Where("group_id = ? AND user_id = ?", removeMemberDTO.GroupID, removeMemberDTO.MemberID).First(&targetMember).Error; err == nil {
			if targetMember.Role == "admin" {
				model.SendResponse(c, http.StatusForbidden, model.Error("管理员不能移除其他管理员"))
				return
			}
		}
	}

	// 移除成员
	if err := h.groupService.RemoveMember(removeMemberDTO.GroupID, removeMemberDTO.MemberID); err != nil {
		config.Logger.Error(err)
		model.SendResponse(c, http.StatusInternalServerError, model.Error(err.Error()))
		return
	}

	model.SendResponse(c, http.StatusOK, model.Success("成员移除成功", nil))
}
