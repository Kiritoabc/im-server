package db

import (
	"time"
)

// Group 群组表结构体
type Group struct {
	ID          uint      `gorm:"primaryKey" json:"id"`                                                                                                                            // 主键
	Name        string    `gorm:"not null" json:"name"`                                                                                                                            // 群组名称，不能为空
	OwnerID     uint      `gorm:"not null" json:"owner_id"`                                                                                                                        // 群主的用户ID，不能为空
	GroupAvatar string    `gorm:"default:'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQ3w2fqb71MsCj97IKLAUXoI6BS4IfeCeEoq_XGS3X2CErGlYyP4xxX4eQ&s'" json:"group_avatar"` // 群组头像，不能为空
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`                                                                                                                // 创建时间
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`                                                                                                                // 更新时间
	//DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`          // 删除时间
}

func CreateGroup(group Group) error {
	return DB.Create(&group).Error
}
