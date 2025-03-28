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
	}
	return s.db.Create(&message).Error
}

// SaveGroupMessage 保存群聊消息到数据库
func (s *MessageService) SaveGroupMessage(senderID, groupID uint, content string) error {
	message := db.Message{
		SenderID: senderID,
		// GroupID:     groupID,
		Content:     content,
		MessageType: "group",
	}
	return s.db.Create(&message).Error
}

// GetGroupMembers 获取群组成员
func (s *MessageService) GetGroupMembers(groupID int) ([]db.GroupMember, error) {
	var members []db.GroupMember
	err := s.db.Where("group_id = ?", groupID).Find(&members).Error
	return members, err
}
