package dto

// DeleteFriendDTO 删除好友请求参数
type DeleteFriendDTO struct {
	FriendID uint `json:"friend_id" binding:"required"`
}

// UpdateFriendGroupDTO 更新好友分组请求参数
type UpdateFriendGroupDTO struct {
	FriendID uint   `json:"friend_id" binding:"required"`
	GroupID  uint   `json:"group_id" binding:"required"`
	Remark   string `json:"remark"`
}
