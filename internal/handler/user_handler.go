package handler

import (
	"github.com/gin-gonic/gin"
	"im-system/internal/config"
	"im-system/internal/model"
	"im-system/internal/model/db"
	"im-system/internal/model/dto"
	"im-system/internal/service"
	"net/http"
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
	var addFriendDTO dto.AddFriendDTO
	if err := c.ShouldBindJSON(&addFriendDTO); err != nil {
		model.SendResponse(c, http.StatusBadRequest, model.Error("无效的请求"))
		return
	}
	if addFriendDTO.GroupID == 0 {
		model.SendResponse(c, http.StatusBadRequest, model.Error("请选择好友分组"))
		return
	}
	userID, exists := c.Get("user_id")
	if !exists {
		model.SendResponse(c, http.StatusUnauthorized, model.Error("用户未登录"))
		return
	}
	addFriendDTO.UserID = userID.(uint)

	if err := h.userService.AddFriend(addFriendDTO); err != nil {
		config.Logger.Error(err)
		model.SendResponse(c, http.StatusInternalServerError, model.Error(err.Error()))
		return
	}

	model.SendResponse(c, http.StatusOK, model.Success("发送好友添加请求", nil))
}
