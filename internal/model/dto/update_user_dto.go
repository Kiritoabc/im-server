package dto

// UpdateUserDTO 修改用户请求参数
type UpdateUserDTO struct {
	AvatarURL string `json:"avatar_url"` // 用户头像URL，允许为空
	Username  string `json:"username"`   // 用户名，允许为空
	Gender    string `json:"gender"`     // 用户性别，允许为空
	Birthday  string `json:"birthday"`   // 用户生日，允许为空
	Bio       string `json:"bio"`        // 用户个人简介，允许为空
	City      string `json:"city"`       // 用户所在城市，允许为空
}
