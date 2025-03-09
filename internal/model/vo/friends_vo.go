package vo

import "time"

type FriendVO struct {
	ID          uint      `json:"id"`           // 好友ID
	Name        string    `json:"name"`         // 好友姓名
	Avatar      string    `json:"avatar"`       // 好友头像URL
	Status      string    `json:"status"`       // 好友状态（在线/离线）
	Email       string    `json:"email"`        // 好友邮箱
	PhoneNumber string    `json:"phone_number"` // 好友电话号码
	Bio         string    `json:"bio"`          // 好友个人简介
	Gender      string    `json:"gender"`       // 好友性别
	City        string    `json:"city"`         // 好友所在城市
	CreatedAt   time.Time `json:"created_at"`   // 创建时间
	UpdatedAt   time.Time `json:"updated_at"`   // 更新时间
}
