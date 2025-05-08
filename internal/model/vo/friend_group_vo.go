package vo

import (
	"im-system/internal/model/db"
	"time"
)

// FriendGroupVO 好友分组视图对象

type FriendGroupVO struct {
	GroupID   uint      `json:"group_id"`   // 分组ID
	GroupName string    `json:"group_name"` // 分组名称
	Members   []db.User `json:"members"`    // 分组中的好友信息
	CreatedAt time.Time `json:"created_at"` // 创建时间
}
