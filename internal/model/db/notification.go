package db

import (
	"gorm.io/gorm"
	"time"
)

// Notification 通知表结构体
type Notification struct {
	ID        uint           `gorm:"primaryKey" json:"id"`                 // 通知ID，自增主键
	UserID    uint           `gorm:"not null" json:"user_id"`              // 用户ID，不能为空
	Type      string         `gorm:"default:'friend_request'" json:"type"` // 通知类型，默认为 'friend_request'
	Content   string         `gorm:"default:''" json:"content"`            // 通知内容，允许为空
	IsRead    bool           `gorm:"default:false" json:"is_read"`         // 是否已读，默认为 false
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`     // 创建时间
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`     // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`              // 删除时间
}

func CreateNotification(notification Notification) error {
	return DB.Create(&notification).Error
}
