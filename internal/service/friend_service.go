package service

import (
	"encoding/json"
	"errors"
	"im-system/internal/config"
	"im-system/internal/model/db"
	"im-system/internal/model/dto"
	"im-system/internal/model/vo"
	"log"

	"gorm.io/gorm"
)

type FriendService struct {
	db *gorm.DB
}

func NewFriendService() *FriendService {
	return &FriendService{
		db: db.DB,
	}
}

type Message struct {
	Id          int    `json:"id"`
	SenderId    int    `json:"senderId"`
	ReceiverID  int    `json:"receiverId"` // 用于私聊
	GroupID     int    `json:"groupId"`    // 用于群聊
	SenderName  string `json:"senderName"`
	Avatar      string `json:"avatar"`
	Content     string `json:"content"`
	MessageType string `json:"messageType"` // "private" 或 "group"
}

// GetUserFriendsGroups 获取用户的好友
func (s *FriendService) GetUserFriendsGroups(userID uint) ([]db.FriendGroup, error) {
	var groups []db.FriendGroup
	if err := s.db.Where("user_id = ?", userID).Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

// GetFriendGroupsWithMembers 获取用户的好友分组及其成员
func (s *FriendService) GetFriendGroupsWithMembers(userID uint) ([]vo.FriendShipGroupVO, error) {
	var resp []vo.FriendShipGroupVO

	var groups []db.FriendGroup
	if err := s.db.Where("user_id = ?", userID).Find(&groups).Error; err != nil {
		return nil, err
	}

	// 查询当前用户的所有好友
	for i := range groups {
		resp = append(resp, vo.FriendShipGroupVO{
			GroupID:   groups[i].ID,
			GroupName: groups[i].GroupName,
			CreatedAt: groups[i].CreatedAt,
			Members:   []vo.FriendMemberVO{},
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
				// 将用户信息添加到分组中，同时包含备注信息
				resp[i].Members = append(resp[i].Members, vo.FriendMemberVO{
					ID:          user.ID,
					Username:    user.Username,
					AvatarURL:   user.AvatarURL,
					Email:       user.Email,
					PhoneNumber: user.PhoneNumber,
					Bio:         user.Bio,
					Gender:      user.Gender,
					City:        user.City,
					Remark:      friend.Remark, // 添加备注信息
					CreatedAt:   user.CreatedAt,
					UpdatedAt:   user.UpdatedAt,
				})
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
			// 查询历史消息，按照时间排序
			var messages []db.Message
			if err := s.db.Where("sender_id =? and receiver_user_id =?", userId, friend.FriendID).
				Or("sender_id =? and receiver_user_id =?", friend.FriendID, userId).
				Order("created_at desc").Find(&messages).Error; err != nil {
				return nil, err
			}
			msgs, err := reverseMessage(messages)
			if err != nil {
				log.Println("reverseMessage error:", err)
				return nil, err
			}
			// 将用户信息添加到分组中
			resp = append(resp, vo.ChatVO{
				ID:          int(user.ID),
				Name:        user.Username,
				Avatar:      user.AvatarURL,
				LastMessage: "",
				Type:        "personal",
				// todo: 查询和好友所有的聊天记录
				Messages: msgs,
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
		// 查询历史消息
		var messages []db.Message
		if err := s.db.Where("receiver_group_id =?", group.GroupID).
			Order("created_at desc").
			Find(&messages).Error; err != nil {
			return nil, err
		}
		msgs, err := reverseMessage(messages)
		if err != nil {
			return nil, err
		}
		resp = append(resp, vo.ChatVO{
			ID:          int(groupInfo.ID),
			Name:        groupInfo.Name,
			Avatar:      groupInfo.GroupAvatar,
			LastMessage: "",
			Type:        "group",
			// todo: 查询和好友所有的聊天记录
			Messages: msgs,
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

// reverseMessage 反转消息
func reverseMessage(messages []db.Message) ([]vo.MessageVO, error) {
	var resp []vo.MessageVO
	for i := len(messages) - 1; i >= 0; i-- {
		var msg Message
		if err := json.Unmarshal([]byte(messages[i].Content), &msg); err != nil {
			log.Println("Unmarshal error:", err)
			return nil, err
		}
		resp = append(resp, vo.MessageVO{
			ID:         msg.Id,
			Content:    msg.Content,
			SenderID:   msg.SenderId,
			SenderName: msg.SenderName,
			Avatar:     msg.Avatar,
			CreatedAt:  messages[i].CreatedAt,
		})
	}
	return resp, nil
}

// SearchFriendGroups 搜索好友分组
func (s *FriendService) SearchFriendGroups(userID uint, keyword string) ([]vo.FriendGroupVO, error) {
	var resp []vo.FriendGroupVO

	// 查询包含关键词的分组
	var groups []db.FriendGroup
	if err := s.db.Where("user_id = ? AND group_name LIKE ?", userID, "%"+keyword+"%").Find(&groups).Error; err != nil {
		return nil, err
	}

	// 查询每个分组的好友
	for i := range groups {
		resp = append(resp, vo.FriendGroupVO{
			GroupID:   groups[i].ID,
			GroupName: groups[i].GroupName,
			CreatedAt: groups[i].CreatedAt,
			Members:   []db.User{},
		})

		// 查询分组中的好友
		var friends []db.Friendship
		if err := s.db.Where("user_id = ? AND group_id = ? AND status = 'accepted'", userID, groups[i].ID).Find(&friends).Error; err != nil {
			return nil, err
		}

		// 查询好友的详细信息
		for _, friend := range friends {
			var user db.User
			if err := s.db.First(&user, friend.FriendID).Error; err == nil {
				resp[i].Members = append(resp[i].Members, user)
			}
		}
	}

	return resp, nil
}

// UpdateFriendGroup 更新好友分组
func (s *FriendService) UpdateFriendGroup(userID uint, dto dto.UpdateFriendGroupDTO) error {
	// 检查好友关系是否存在
	var friendship db.Friendship
	if err := s.db.Where("user_id = ? AND friend_id = ?", userID, dto.FriendID).First(&friendship).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("好友关系不存在")
		}
		return err
	}

	// 检查目标分组是否存在
	var group db.FriendGroup
	if err := s.db.Where("id = ? AND user_id = ?", dto.GroupID, userID).First(&group).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("分组不存在")
		}
		return err
	}

	// 更新好友分组和备注
	updates := map[string]interface{}{
		"group_id": dto.GroupID,
	}
	if dto.Remark != "" {
		updates["remark"] = dto.Remark
	}

	if err := s.db.Model(&friendship).Updates(updates).Error; err != nil {
		config.Logger.Error(err)
		return errors.New("更新好友分组失败")
	}

	return nil
}
