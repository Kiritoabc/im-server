package service

import (
	"im-system/internal/model/db"

	"gorm.io/gorm"
)

type MessageService struct {
	db *gorm.DB
}

func NewMessageService() *MessageService {
	return &MessageService{db: db.DB}
}

// SaveMessage 保存消息到数据库
func (s *MessageService) SaveMessage(senderID, receiverID uint, content string) error {
	message := db.Message{
		SenderID:       senderID,
		ReceiverUserID: receiverID,
		Content:        content,
		MessageType:    "text",
	}
	return s.db.Create(&message).Error
}

// SaveGroupMessage 保存群聊消息到数据库
func (s *MessageService) SaveGroupMessage(senderID, groupID uint, content string) error {
	message := db.Message{
		SenderID:        senderID,
		ReceiverGroupID: groupID,
		Content:         content,
		MessageType:     "text",
	}
	return s.db.Create(&message).Error
}

// GetGroupMembers 获取群组成员
func (s *MessageService) GetGroupMembers(groupID int) ([]db.GroupMember, error) {
	var members []db.GroupMember
	err := s.db.Where("group_id = ?", groupID).Find(&members).Error
	return members, err
}

// GetChatMessages 获取聊天记录
func (s *MessageService) GetChatMessages(chatType string, userID, toID uint) ([]db.Message, error) {
	var messages []db.Message
	query := s.db

	if chatType == "private" {
		// 私聊消息：获取双向的消息记录
		query = query.Where(
			"(sender_id = ? AND receiver_user_id = ?) OR (sender_id = ? AND receiver_user_id = ?)",
			userID, toID, toID, userID,
		).Where("receiver_group_id IS NULL")
	} else if chatType == "group" {
		// 群聊消息：获取指定群的消息记录
		query = query.Where("receiver_group_id = ?", toID)
	}

	// 按时间倒序排序并限制最近100条消息
	err := query.Order("created_at DESC").
		Limit(100).
		Find(&messages).Error

	if err != nil {
		return nil, err
	}

	return messages, nil
}

// GetMessageWithSenderInfo 获取带发送者信息的聊天记录
func (s *MessageService) GetMessageWithSenderInfo(chatType string, userID, toID uint) ([]db.Message, error) {

	var messages []db.Message
	if chatType == "private" {
		// 查询历史消息，按照时间排序
		if err := s.db.Where("sender_id =? and receiver_user_id =?", userID, toID).
			Or("sender_id =? and receiver_user_id =?", toID, userID).
			Order("created_at desc").Find(&messages).Error; err != nil {
			return nil, err
		}
	} else if chatType == "group" {
		if err := s.db.Where("receiver_group_id =?", toID).
			Order("created_at desc").
			Find(&messages).Error; err != nil {
			return nil, err
		}
	}

	return messages, nil
}
