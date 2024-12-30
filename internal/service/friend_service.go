package service

import (
	"gorm.io/gorm"
	"im-system/internal/model/db"
	"im-system/internal/model/vo"
)

type FriendService struct {
	db *gorm.DB
}

func NewFriendService() *FriendService {
	return &FriendService{
		db: db.DB,
	}
}

// GetFriendGroupsWithMembers 获取用户的好友分组及其成员
func (s *FriendService) GetFriendGroupsWithMembers(userID uint) ([]vo.FriendGroupVO, error) {
	var resp []vo.FriendGroupVO

	var groups []db.FriendGroup
	if err := s.db.Where("user_id = ?", userID).Find(&groups).Error; err != nil {
		return nil, err
	}

	// 查询当前用户的所有好友
	for i := range groups {
		resp = append(resp, vo.FriendGroupVO{
			GroupID:   groups[i].ID,
			GroupName: groups[i].GroupName,
			CreatedAt: groups[i].CreatedAt,
			Members:   []db.User{},
		})
		var friends []db.Friendship
		// 查询所有接受的好友
		if err := s.db.Where("group_id = ? and status = 'accepted'", groups[i].ID).Find(&friends).Error; err != nil {
			return nil, err
		}
		// 这里可以进一步查询每个好友的基本信息
		for _, friend := range friends {
			var user db.User
			if err := s.db.First(&user, friend.FriendID).Error; err == nil {
				// 将用户信息添加到分组中
				resp[i].Members = append(resp[i].Members, user)
			}
		}
	}

	return resp, nil
}
