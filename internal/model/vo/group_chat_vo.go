package vo

import (
	"time"
)

// GroupChatVO 群聊视图对象
type GroupChatVO struct {
	ID          uint      `json:"id"`           // 群组ID
	Name        string    `json:"name"`         // 群组名称
	OwnerID     uint      `json:"owner_id"`     // 群主ID
	Role        string    `json:"role"`         // 用户在群组中的角色（Owner/admin/member）
	GroupAvatar string    `json:"group_avatar"` // 群组头像
	CreatedAt   time.Time `json:"created_at"`   // 创建时间
	UpdatedAt   time.Time `json:"updated_at"`   // 更新时间
}
