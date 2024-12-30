package vo

import (
	"im-system/internal/model/db"
	"time"
)

// NotificationVO 通知视图对象
type NotificationVO struct {
	ID        uint      `json:"id"`         // 通知ID
	UserID    uint      `json:"user_id"`    // 用户ID
	Type      string    `json:"type"`       // 通知类型
	Content   string    `json:"content"`    // 通知内容
	IsRead    bool      `json:"is_read"`    // 是否已读
	CreatedAt time.Time `json:"created_at"` // 创建时间
	Sender    db.User   `json:"sender"`     // 发送请求的用户信息
}
