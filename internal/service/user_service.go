package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"im-system/internal/config"
	"im-system/internal/middle"
	"im-system/internal/model/db"
	"im-system/internal/model/dto"
	"im-system/internal/utils"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService() *UserService {
	return &UserService{
		db: db.DB,
	}
}

func (s *UserService) Register(userInfo db.User) error {
	var existUser db.User
	if err := s.db.Where("phone_number = ?", userInfo.PhoneNumber).First(&existUser).Error; err == nil {
		return errors.New("手机号已经被注册")
	}

	userInfo.PasswordHash = utils.HashPassword(userInfo.PasswordHash)

	// 防止邮箱为空
	if userInfo.Email == "" {
		userInfo.Email = fmt.Sprintf("%s@imSystem.com", userInfo.Username)
	}

	// 默认为 male
	if userInfo.Gender == "" {
		userInfo.Gender = "male"
	}
	if userInfo.AvatarURL == "" {
		userInfo.AvatarURL = "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQ3w2fqb71MsCj97IKLAUXoI6BS4IfeCeEoq_XGS3X2CErGlYyP4xxX4eQ&s"
	}

	// 注册用户
	if err := s.db.Create(&userInfo).Error; err != nil {
		config.Logger.Error(err)
		return err
	}

	// 创建默认的好友分组
	defaultGroup := db.FriendGroup{
		UserID:    userInfo.ID, // 使用新注册用户的ID
		GroupName: "我的好友",      // 默认分组名称
	}

	if err := s.db.Create(&defaultGroup).Error; err != nil {
		config.Logger.Error(err)
		return err
	}

	return nil
}

func (s *UserService) Login(phoneNumber, password string) (uint, string, error) {
	var user db.User
	if err := s.db.Where("phone_number = ?", phoneNumber).First(&user).Error; err != nil {
		return 0, "", errors.New("用户不存在")
	}

	if !utils.ComparePassword(user.PasswordHash, password) {
		return 0, "", errors.New("密码错误")
	}

	// 生成 JWT
	token, err := middle.GenerateJWT(user.ID)
	if err != nil {
		return 0, "", err
	}

	// 将 token 和用户信息存储到 Redis
	if err := middle.SetTokenToRedis(user.ID, user); err != nil {
		return 0, "", err
	}

	// 返回用户ID和token
	return user.ID, token, nil
}

// GetUserInfo 获取用户信息
func (s *UserService) GetUserInfo(userID uint) (db.User, error) {
	var user db.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return user, errors.New("用户不存在")
	}
	return user, nil
}

// AddFriend 添加好友,todo: 事务完善
func (s *UserService) AddFriend(addFriendDto dto.AddFriendDTO) error {
	var friend db.User
	// 1.查询好友是否存在
	if err := s.db.First(&friend, addFriendDto.FriendID).Error; err != nil {
		return errors.New("好友不存在")
	}

	// 查询分组是否存在
	var group db.FriendGroup
	if err := s.db.First(&group, addFriendDto.GroupID).Error; err != nil {
		return errors.New("分组不存在")
	}

	// 2.创建好友关系，不创建关系表，直接创建通知
	friendship := db.Friendship{
		UserID:   addFriendDto.UserID,
		FriendID: addFriendDto.FriendID,
		Remark:   addFriendDto.Remark,
		GroupID:  addFriendDto.GroupID,
		Status:   "pending", // 发送中
	}

	if err := s.db.Create(&friendship).Error; err != nil {
		return err
	}

	// 3.创建通知。发送通知
	notification := db.Notification{
		SenderID:   addFriendDto.UserID,
		ReceiverID: addFriendDto.FriendID,
		Type:       "friend_request",
		Content:    addFriendDto.Content,
	}

	if err := s.db.Create(&notification).Error; err != nil {
		config.Logger.Error(err)
		return err
	}

	return nil
}

// Logout 处理用户退出登录
func (s *UserService) Logout(c *gin.Context, userId uint) error {
	// 删除 Redis 中的用户信息
	if err := config.RedisClient.Del(c, middle.GetRedisUserInfoKey(userId)).Err(); err != nil {
		config.Logger.Error(err)
		return err
	}
	return nil
}
