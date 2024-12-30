package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"im-system/internal/config"
	"im-system/internal/model"
	"im-system/internal/model/db"
	"im-system/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(c *gin.Context) {
	var user db.User
	if err := c.ShouldBindJSON(&user); err != nil {
		model.SendResponse(c, http.StatusBadRequest, model.Error("无效的请求"))
		return
	}

	if err := h.userService.Register(user); err != nil {
		config.Logger.Error(err)
		model.SendResponse(c, http.StatusInternalServerError, model.Error(err.Error()))
		return
	}

	model.SendResponse(c, http.StatusOK, model.Success("注册成功", nil))
}

func (h *UserHandler) Login(c *gin.Context) {
	var loginData struct {
		PhoneNumber string `json:"phone_number"`
		Password    string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		model.SendResponse(c, http.StatusBadRequest, model.Error("无效的请求"))
		return
	}

	userID, token, err := h.userService.Login(loginData.PhoneNumber, loginData.Password)
	if err != nil {
		config.Logger.Error(err)
		model.SendResponse(c, http.StatusUnauthorized, model.Error(err.Error()))
		return
	}

	model.SendResponse(c, http.StatusOK, model.Success("登录成功", gin.H{"user_id": userID, "token": token}))
}

// GetUserInfo 获取用户信息
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		model.SendResponse(c, http.StatusUnauthorized, model.Error("未提供用户ID"))
		return
	}

	user, err := h.userService.GetUserInfo(userID.(uint))
	if err != nil {
		config.Logger.Error(err)
		model.SendResponse(c, http.StatusNotFound, model.Error(err.Error()))
		return
	}

	model.SendResponse(c, http.StatusOK, model.Success("获取用户信息成功", user))
}

// AddFriend 添加好友
func (h *UserHandler) AddFriend(c *gin.Context) {
	userIDStr := c.Param("user_id")
	friendIDStr := c.Param("friend_id")
	remark := c.Query("remark")       // 从查询参数获取备注
	groupIDStr := c.Query("group_id") // 从查询参数获取分组ID。这个前端需要传过来

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		model.SendResponse(c, http.StatusBadRequest, model.Error("无效的用户ID"))
		return
	}

	friendID, err := strconv.ParseUint(friendIDStr, 10, 32)
	if err != nil {
		model.SendResponse(c, http.StatusBadRequest, model.Error("无效的好友ID"))
		return
	}

	var groupID *uint
	if groupIDStr != "" {
		id, err := strconv.ParseUint(groupIDStr, 10, 32)
		if err == nil {
			groupID = new(uint)
			*groupID = uint(id)
		}
	}

	if err := h.userService.AddFriend(uint(userID), uint(friendID), remark, groupID); err != nil {
		config.Logger.Error(err)
		model.SendResponse(c, http.StatusInternalServerError, model.Error(err.Error()))
		return
	}

	model.SendResponse(c, http.StatusOK, model.Success("好友添加成功", nil))
}
