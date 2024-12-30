package dto

// AddFriendDTO 添加好友请求参数
type AddFriendDTO struct {
	UserID   uint   `json:"user_id"`   // 用户ID
	FriendID uint   `json:"friend_id"` // 好友ID
	Remark   string `json:"remark"`    // 备注
	GroupID  uint   `json:"group_id"`  // 分组ID
	Content  string `json:"content"`   // 消息内容
}
