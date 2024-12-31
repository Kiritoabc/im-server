package service

import (
	"gorm.io/gorm"
	"im-system/internal/model/db"
)

type FriendGroupService struct {
	db *gorm.DB
}

func NewFriendShipService() *FriendGroupService {
	return &FriendGroupService{
		db: db.DB,
	}
}

// CreateFriendGroup 创建好友分组
func (s *FriendGroupService) CreateFriendGroup(userID uint, groupName string) error {
	group := db.FriendGroup{
		UserID:    userID,
		GroupName: groupName,
	}

	return s.db.Create(&group).Error
}

// GetFriendGroupsWithMembers 获取好友的所有分组
func (s *FriendGroupService) GetFriendGroupsWithMembers(userId uint) ([]db.FriendGroup, error) {
	resp := make([]db.FriendGroup, 0)
	if err := s.db.Where("user_id = ?", userId).Find(&resp).Error; err != nil {
		return resp, err
	}
	return resp, nil
}
