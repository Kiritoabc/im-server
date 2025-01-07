package vo

import (
	"time"
)

// GroupVO 群组视图对象
type GroupVO struct {
	ID          uint      `json:"id"`           // 群组ID
	Name        string    `json:"name"`         // 群组名称
	OwnerID     uint      `json:"owner_id"`     // 群主ID
	GroupAvatar string    `json:"group_avatar"` // 群组头像
	CreatedAt   time.Time `json:"created_at"`   // 创建时间
	UpdatedAt   time.Time `json:"updated_at"`   // 更新时间
}
