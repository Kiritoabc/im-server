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
		// 查询所有接受的好友(我的好友)
		if err := s.db.Where("user_id = ? and group_id = ? and status = 'accepted'", userID, groups[i].ID).Find(&friends).Error; err != nil {
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

// GetUserFriendsChat 获取用户的好友聊天
func (s *FriendService) GetUserFriendsChat(userId uint) ([]vo.ChatVO, error) {
	var resp []vo.ChatVO

	// 查询所有接受的好友
	var friends []db.Friendship
	if err := s.db.Where("user_id =? and status = 'accepted'", userId).Find(&friends).Error; err != nil {
		return nil, err
	}
	// 这里可以进一步查询每个好友的基本信息
	for _, friend := range friends {
		var user db.User
		if err := s.db.First(&user, friend.FriendID).Error; err == nil {
			// 将用户信息添加到分组中
			resp = append(resp, vo.ChatVO{
				ID:          int(user.ID),
				Name:        user.Username,
				Avatar:      user.AvatarURL,
				LastMessage: "",
				Type:        "personal",
				// todo: 查询和好友所有的聊天记录
				Messages: []vo.MessageVO{},
			})
		}
	}
	// 查询我所在的群聊
	var groupMembers []db.GroupMember
	if err := s.db.Where("user_id =?", userId).Find(&groupMembers).Error; err != nil {
		return nil, err
	}
	// 这里可以进一步查询每个好友的基本信息
	for _, group := range groupMembers {
		groupInfo := db.Group{}
		if err := s.db.First(&groupInfo, group.GroupID).Error; err != nil {
			return nil, err
		}
		// 将用户信息添加到分组中
		resp = append(resp, vo.ChatVO{
			ID:          int(groupInfo.ID),
			Name:        groupInfo.Name,
			Avatar:      groupInfo.GroupAvatar,
			LastMessage: "",
			Type:        "group",
			// todo: 查询和好友所有的聊天记录
			Messages: []vo.MessageVO{},
		})
	}

	return resp, nil
}

// GetUserFriends 获取用户的好友
func (s *FriendService) GetUserFriends(userID uint) ([]vo.FriendVO, error) {
	var resp []vo.FriendVO

	// 查询所有接受的好友
	var friends []db.Friendship
	if err := s.db.Where("user_id = ? and status = 'accepted'", userID).Find(&friends).Error; err != nil {
		return nil, err
	}

	// 查询好友信息
	for _, friend := range friends {
		var user db.User
		if err := s.db.First(&user, friend.FriendID).Error; err == nil {
			resp = append(resp, vo.FriendVO{
				ID:          user.ID,
				Name:        user.Username,
				Avatar:      user.AvatarURL,
				Status:      "离线", // todo: 这里需要根据实际情况获取用户状态
				Email:       user.Email,
				PhoneNumber: user.PhoneNumber,
				Bio:         user.Bio,
				Gender:      user.Gender,
				City:        user.City,
				CreatedAt:   user.CreatedAt,
				UpdatedAt:   user.UpdatedAt,
			})
		}
	}

	return resp, nil
}
