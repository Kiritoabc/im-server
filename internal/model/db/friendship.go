package db

import (
	"gorm.io/gorm"
	"time"
)

// Friendship 好友关系表结构体
type Friendship struct {
	ID        uint           `gorm:"primaryKey" json:"id"`             // 好友关系ID，自增主键
	UserID    uint           `gorm:"not null" json:"user_id"`          // 用户ID，不能为空
	FriendID  uint           `gorm:"not null" json:"friend_id"`        // 好友的用户ID，不能为空
	Status    string         `gorm:"default:'pending'" json:"status"`  // 好友关系状态，默认为 'pending'
	GroupID   uint           `gorm:"default:NULL" json:"group_id"`     // 分组ID，允许为空
	Remark    string         `gorm:"default:''" json:"remark"`         // 好友备注，允许为空
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"` // 创建时间
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"` // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`          // 删除时间
}

func CreateFriendship(friendship Friendship) error {
	return DB.Create(&friendship).Error
}
