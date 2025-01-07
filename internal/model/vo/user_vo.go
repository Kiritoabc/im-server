package vo

import (
	"time"
)

// UserVO 用户视图对象
type UserVO struct {
	ID          uint      `json:"id"`           // 用户ID
	Username    string    `json:"username"`     // 用户名
	Email       string    `json:"email"`        // 用户邮箱
	PhoneNumber string    `json:"phone_number"` // 用户电话号码
	AvatarURL   string    `json:"avatar_url"`   // 用户头像URL
	Bio         string    `json:"bio"`          // 用户个人简介
	Gender      string    `json:"gender"`       // 用户性别
	CreatedAt   time.Time `json:"created_at"`   // 创建时间
	UpdatedAt   time.Time `json:"updated_at"`   // 更新时间
	City        string    `json:"city"`         // 用户城市
}
