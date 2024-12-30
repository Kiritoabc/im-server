package db

import (
	"time"
)

// Notification 通知表结构体
type Notification struct {
	ID         uint      `gorm:"primaryKey" json:"id"`                                                                                       // 通知ID，自增主键
	SenderID   uint      `gorm:"not null" json:"sender_id"`                                                                                  // 发送者ID，不能为空
	ReceiverID uint      `gorm:"not null" json:"receiver_id"`                                                                                // 接收者ID，不能为空
	Type       string    `gorm:"type:enum('message', 'friend_request', 'group_request', 'rejected', 'other');default:'message'" json:"type"` // 通知类型，默认为 'message',rejected
	Content    string    `gorm:"default:''" json:"content"`                                                                                  // 通知内容，允许为空
	IsRead     bool      `gorm:"default:false" json:"is_read"`                                                                               // 是否已读，默认为 false
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`                                                                           // 创建时间
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`                                                                           // 更新时间
}

// CreateNotification 创建通知
func CreateNotification(notification Notification) error {
	return DB.Create(&notification).Error
}
