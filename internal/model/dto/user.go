package dto

// DeleteFriendDTO 删除好友请求参数
type DeleteFriendDTO struct {
	FriendID uint `json:"friend_id" binding:"required"`
}
