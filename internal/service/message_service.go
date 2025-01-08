package service

import (
	"gorm.io/gorm"
	"im-system/internal/model/db"
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
