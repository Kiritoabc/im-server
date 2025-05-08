package vo

import "time"

// FriendShipGroupVO 好友分组视图对象
type FriendShipGroupVO struct {
	GroupID   uint             `json:"group_id"`
	GroupName string           `json:"group_name"`
	CreatedAt time.Time        `json:"created_at"`
	Members   []FriendMemberVO `json:"members"`
}

// FriendMemberVO 好友成员视图对象
type FriendMemberVO struct {
	ID          uint      `json:"id"`
	Username    string    `json:"username"`
	AvatarURL   string    `json:"avatar_url"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	Bio         string    `json:"bio"`
	Gender      string    `json:"gender"`
	City        string    `json:"city"`
	Remark      string    `json:"remark"` // 好友备注
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
