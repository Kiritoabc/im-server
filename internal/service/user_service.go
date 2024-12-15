package service

import (
	"context"
	"errors"
	"time"

	"im-system/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	// 在这里可以注入其他依赖
}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) Register(username, password string) error {
	var existUser model.User
	if err := model.DB.Where("username = ?", username).First(&existUser).Error; err == nil {
		return errors.New("用户名已存在")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := model.User{
		Username: username,
		Password: string(hashedPassword),
	}

	return model.DB.Create(&user).Error
}

func (s *UserService) Login(username, password string) (uint, error) {
	var user model.User
	if err := model.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return 0, errors.New("用户不存在")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return 0, errors.New("密码错误")
	}

	return user.ID, nil
}

func (s *UserService) AddFriend(userID, friendID uint) error {
	var friend model.User
	if err := model.DB.First(&friend, friendID).Error; err != nil {
		return errors.New("好友不存在")
	}

	friendship := model.Friend{
		UserID:   userID,
		FriendID: friendID,
	}

	return model.DB.Create(&friendship).Error
}

func (s *UserService) SendMessage(fromUserID, toUserID uint, content string) error {
	message := model.Message{
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Content:    content,
		CreatedAt:  time.Now(),
	}

	_, err := model.MongoDB.Collection("messages").InsertOne(context.Background(), message)
	return err
}