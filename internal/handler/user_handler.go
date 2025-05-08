package handler

import (
	"fmt"
	"im-system/internal/config"
	"im-system/internal/model"
	"im-system/internal/model/db"
	"im-system/internal/model/dto"
	"im-system/internal/service"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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

// Logout 处理用户退出登录
func (h *UserHandler) Logout(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		model.SendResponse(c, http.StatusUnauthorized, model.Error("用户未登录"))
		return
	}

	if err := h.userService.Logout(c, userID.(uint)); err != nil {
		config.Logger.Error(err)
		model.SendResponse(c, http.StatusInternalServerError, model.Error("退出登录失败"))
		return
	}

	model.SendResponse(c, http.StatusOK, model.Success("退出登录成功", nil))
}

// UpdateUserInfo 更新用户信息
func (h *UserHandler) UpdateUserInfo(c *gin.Context) {
	var updateUserDTO dto.UpdateUserDTO
	if err := c.ShouldBindJSON(&updateUserDTO); err != nil {
		model.SendResponse(c, http.StatusBadRequest, model.Error("无效的请求"))
		return
	}

	// 从上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		model.SendResponse(c, http.StatusUnauthorized, model.Error("用户未登录"))
		return
	}

	user := db.User{
		ID:          userID.(uint),
		Username:    updateUserDTO.Username,
		AvatarURL:   updateUserDTO.AvatarURL,
		DateOfBirth: updateUserDTO.Birthday,
		City:        updateUserDTO.City,
		Bio:         updateUserDTO.Bio,
		Gender:      updateUserDTO.Gender,
	}

	if err := h.userService.UpdateUserInfo(userID.(uint), user); err != nil {
		config.Logger.Error(err)
		model.SendResponse(c, http.StatusInternalServerError, model.Error(err.Error()))
		return
	}

	model.SendResponse(c, http.StatusOK, model.Success("用户信息更新成功", nil))
}

// QueryUserAndGroup 查询用户和群聊信息
func (h *UserHandler) QueryUserAndGroup(c *gin.Context) {
	var queryDTO dto.QueryUserAndGroupDTO
	if err := c.ShouldBindJSON(&queryDTO); err != nil {
		model.SendResponse(c, http.StatusBadRequest, model.Error("无效的请求"))
		return
	}

	result, err := h.userService.QueryUserAndGroup(queryDTO)
	if err != nil {
		config.Logger.Error(err)
		model.SendResponse(c, http.StatusInternalServerError, model.Error(err.Error()))
		return
	}

	model.SendResponse(c, http.StatusOK, model.Success("查询成功", result))
}

// DeleteFriend 删除好友
func (h *UserHandler) DeleteFriend(c *gin.Context) {
	var deleteFriendDTO dto.DeleteFriendDTO
	if err := c.ShouldBindJSON(&deleteFriendDTO); err != nil {
		model.SendResponse(c, http.StatusBadRequest, model.Error("无效的请求"))
		return
	}

	// 从上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		model.SendResponse(c, http.StatusUnauthorized, model.Error("用户未登录"))
		return
	}

	if err := h.userService.DeleteFriend(userID.(uint), deleteFriendDTO.FriendID); err != nil {
		config.Logger.Error(err)
		model.SendResponse(c, http.StatusInternalServerError, model.Error(err.Error()))
		return
	}

	model.SendResponse(c, http.StatusOK, model.Success("删除好友成功", nil))
}

// UploadAvatar 上传用户头像
func (h *UserHandler) UploadAvatar(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		model.SendResponse(c, http.StatusUnauthorized, model.Error("用户未登录"))
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("avatar")
	if err != nil {
		model.SendResponse(c, http.StatusBadRequest, model.Error("请选择要上传的头像文件"))
		return
	}

	// 检查文件类型
	if !isImageFile(file.Filename) {
		model.SendResponse(c, http.StatusBadRequest, model.Error("只支持上传图片文件"))
		return
	}

	// 生成文件名：用户ID_时间戳.扩展名
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d_%d%s", userID, time.Now().Unix(), ext)

	// 确保 avatars 目录存在
	avatarDir := "static/avatars"
	if err := os.MkdirAll(avatarDir, 0755); err != nil {
		config.Logger.Error(err)
		model.SendResponse(c, http.StatusInternalServerError, model.Error("创建头像目录失败"))
		return
	}

	// 保存文件
	filepath := filepath.Join(avatarDir, filename)
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		config.Logger.Error(err)
		model.SendResponse(c, http.StatusInternalServerError, model.Error("保存头像文件失败"))
		return
	}

	// 更新用户头像URL
	avatarURL := fmt.Sprintf("/static/avatars/%s", filename)
	if err := h.userService.UpdateUserAvatar(userID.(uint), avatarURL); err != nil {
		config.Logger.Error(err)
		model.SendResponse(c, http.StatusInternalServerError, model.Error("更新用户头像失败"))
		return
	}

	model.SendResponse(c, http.StatusOK, model.Success("上传头像成功", gin.H{
		"avatar_url": "http://localhost:8080" + avatarURL,
	}))
}

// isImageFile 检查文件是否为图片
func isImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif"
}
