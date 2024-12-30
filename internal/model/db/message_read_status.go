package db

import (
	"gorm.io/gorm"
	"time"
)

// MessageReadStatus 消息已读状态表结构体
type MessageReadStatus struct {
	ID        uint           `gorm:"primaryKey" json:"id"`             // 主键
	MessageID uint           `gorm:"not null" json:"message_id"`       // 消息ID，不能为空
	UserID    uint           `gorm:"not null" json:"user_id"`          // 用户ID，不能为空
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"` // 创建时间
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"` // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`          // 删除时间
}

func MarkMessageAsRead(status MessageReadStatus) error {
	return DB.Create(&status).Error
}
