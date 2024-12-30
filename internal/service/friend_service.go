package service

import (
	"gorm.io/gorm"
	"im-system/internal/model/db"
)

type FriendService struct {
	db *gorm.DB
}

func NewFriendService() *FriendService {
	return &FriendService{
		db: db.DB,
	}
}

// CreateFriendGroup 创建好友分组
func (s *FriendService) CreateFriendGroup(userID uint, groupName string) error {
	group := db.FriendGroup{
		UserID:    userID,
		GroupName: groupName,
	}

	return s.db.Create(&group).Error
}
