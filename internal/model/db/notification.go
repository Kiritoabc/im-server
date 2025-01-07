package db

import (
	"time"
)

// Notification 代表通知的结构体
type Notification struct {
	ID         uint      `gorm:"primaryKey" json:"id"`                                                                                       // 通知ID，自增主键
	SenderID   uint      `gorm:"not null" json:"sender_id"`                                                                                  // 发送者ID，不能为空
	ReceiverID uint      `gorm:"not null" json:"receiver_id"`                                                                                // 接收者ID，不能为空
	Type       string    `gorm:"type:enum('message', 'friend_request', 'group_request', 'rejected', 'other');default:'message'" json:"type"` // 通知类型，默认为 'message',rejected
	Content    string    `gorm:"default:''" json:"content"`                                                                                  // 通知内容，允许为空
	IsRead     bool      `gorm:"default:false" json:"is_read"`                                                                               // 是否已读，默认为 false
	Status     string    `gorm:"type:enum('pending', 'accepted', 'rejected');default:'pending'" json:"status"`                               // 通知状态，默认为'pending'
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`                                                                           // 创建时间
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`                                                                           // 更新时间
}

// CreateNotification 创建新的通知
// @param notification 通知对象
// @return error 如果创建失败，返回错误信息
func CreateNotification(notification Notification) error {
	return DB.Create(&notification).Error
}
