package dto

// QueryGroupDTO 查询群组请求参数
type QueryGroupDTO struct {
	GroupID   string `json:"group_id"`   // 群组ID，允许为空
	GroupName string `json:"group_name"` // 群组名称，允许为空
}
