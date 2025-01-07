package vo

// UserAndGroupVO 用户和群组视图对象
type UserAndGroupVO struct {
	Users  []UserVO  `json:"users"`  // 用户信息列表
	Groups []GroupVO `json:"groups"` // 群组信息列表
}
