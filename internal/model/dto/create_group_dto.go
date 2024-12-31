package dto

// CreateGroupDTO 创建群组请求参数
type CreateGroupDTO struct {
	Name        string `json:"name"`         // 群组名称
	OwnerID     uint   `json:"owner_id"`     // 群主的用户ID
	GroupAvatar string `json:"group_avatar"` // 群组头像
}
