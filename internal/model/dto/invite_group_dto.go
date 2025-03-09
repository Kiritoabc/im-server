package dto

type InviteGroupDTO struct {
	GroupID   uint   `json:"group_id"`   // 群组ID
	FriendIDs []uint `json:"friend_ids"` // 要邀请的好友ID列表
}
