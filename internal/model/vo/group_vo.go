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

// GroupStats 群组统计信息
type GroupStats struct {
	Active int `json:"active"` // 活跃人数
	Male   int `json:"male"`   // 男性人数
	Local  int `json:"local"`  // 本地人数
}

// PreviewMember 预览成员信息
type PreviewMember struct {
	ID     uint   `json:"id"`     // 成员ID
	Name   string `json:"name"`   // 成员名称
	Avatar string `json:"avatar"` // 成员头像
	Role   string `json:"role"`   // 成员角色
}

// GroupChatDetail 群聊详细信息
type GroupChatDetail struct {
	ID             uint            `json:"id"`             // 群聊ID
	GroupID        string          `json:"groupId"`        // 群组ID
	Name           string          `json:"name"`           // 群组名称
	Avatar         string          `json:"avatar"`         // 群组头像
	MemberCount    int             `json:"memberCount"`    // 成员数量
	Category       string          `json:"category"`       // 群组分类
	Announcement   string          `json:"announcement"`   // 群公告
	Description    string          `json:"description"`    // 群描述
	Stats          GroupStats      `json:"stats"`          // 群组统计信息
	PreviewMembers []PreviewMember `json:"previewMembers"` // 预览成员列表
}

// GroupChatList 群聊列表
type GroupChatList struct {
	CreatedGroups []GroupChatDetail `json:"created_groups"` // 创建的群聊
	ManagedGroups []GroupChatDetail `json:"managed_groups"` // 管理的群聊
	JoinedGroups  []GroupChatDetail `json:"joined_groups"`  // 加入的群聊
}
