package db

import (
	"time"
)

// Message 消息表结构体
type Message struct {
	ID              uint      `gorm:"primaryKey" json:"id"`                  // 主键
	SenderID        uint      `gorm:"not null" json:"sender_id"`             // 发送者ID，不能为空
	ReceiverUserID  uint      `gorm:"default:NULL" json:"receiver_user_id"`  // 接收者的用户ID，允许为空
	ReceiverGroupID uint      `gorm:"default:NULL" json:"receiver_group_id"` // 接收者的群组ID，允许为空
	Content         string    `gorm:"type:text" json:"content"`              // 消息内容，不能为空
	MessageType     string    `gorm:"default:'text'" json:"message_type"`    // 消息类型，默认为 'text'
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`      // 创建时间
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`      // 更新时间
	//DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at"`               // 删除时间
}

func SendMessage(message Message) error {
	return DB.Create(&message).Error
}
