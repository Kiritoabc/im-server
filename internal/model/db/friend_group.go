package db

import (
	"time"
)

// FriendGroup 分组表结构体
type FriendGroup struct {
	ID        uint      `gorm:"primaryKey" json:"id"`             // 主键
	UserID    uint      `gorm:"not null" json:"user_id"`          // 用户ID，不能为空
	GroupName string    `gorm:"not null" json:"group_name"`       // 分组名称，不能为空
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"` // 创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"` // 更新时间
	//DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"` // 删除时间
}

func CreateFriendGroup(group FriendGroup) error {
	return DB.Create(&group).Error
}
