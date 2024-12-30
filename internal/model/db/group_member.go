package db

import (
	"time"
)

// GroupMember 群组成员表结构体
type GroupMember struct {
	ID        uint      `gorm:"primaryKey" json:"id"`             // 主键
	GroupID   uint      `gorm:"not null" json:"group_id"`         // 群组ID，不能为空
	UserID    uint      `gorm:"not null" json:"user_id"`          // 用户ID，不能为空
	Role      string    `gorm:"default:'member'" json:"role"`     // 成员角色，默认为 'member'
	Title     string    `gorm:"default:''" json:"title"`          // 成员称号，允许为空
	Level     int       `gorm:"default:1" json:"level"`           // 成员等级，默认为 1
	Nickname  string    `gorm:"default:''" json:"nickname"`       // 用户在群组中的昵称，允许为空
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"` // 创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"` // 更新时间
	//DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`          // 删除时间
}

func AddGroupMember(member GroupMember) error {
	return DB.Create(&member).Error
}
