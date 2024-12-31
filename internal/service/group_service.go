package service

import (
	"gorm.io/gorm"
	"im-system/internal/model/db"
	"im-system/internal/model/vo"
)

// GroupService 群组服务
type GroupService struct {
	db *gorm.DB
}

// NewGroupService 创建新的群组服务
func NewGroupService() *GroupService {
	return &GroupService{
		db: db.DB,
	}
}

const (
	Admin  = "admin"
	Owner  = "owner"
	Member = "member"
)

var GroupNameMap = map[string]string{
	Owner:  "我创建的群聊",
	Admin:  "我管理的群聊",
	Member: "我加入的群聊",
}

// CreateGroup 创建群组
func (s *GroupService) CreateGroup(group db.Group, userInfo db.User) error {
	// 创建群组记录
	if err := s.db.Create(&group).Error; err != nil {
		return err
	}

	// 创建群组成员记录
	groupMember := db.GroupMember{
		GroupID:  group.ID,
		UserID:   group.OwnerID,
		Nickname: userInfo.Username, // 默认昵称设置为群组名称
		Role:     Owner,             // 群主角色
	}

	return s.db.Create(&groupMember).Error
}

// QueryGroups 查询群组信息
func (s *GroupService) QueryGroups(groupID, groupName string) ([]vo.GroupVO, error) {
	var groups []db.Group
	var query *gorm.DB

	// 构建查询条件
	query = s.db.Model(&db.Group{})
	if groupID != "" {
		query = query.Where("id LIKE ?", "%"+groupID+"%")
	}
	if groupName != "" {
		query = query.Where("name LIKE ?", "%"+groupName+"%")
	}

	if err := query.Find(&groups).Error; err != nil {
		return nil, err
	}

	// 转换为 GroupVO
	var groupVOs []vo.GroupVO
	for _, group := range groups {
		groupVOs = append(groupVOs, vo.GroupVO{
			ID:        group.ID,
			Name:      group.Name,
			OwnerID:   group.OwnerID,
			CreatedAt: group.CreatedAt,
			UpdatedAt: group.UpdatedAt,
		})
	}

	return groupVOs, nil
}

// GetUserGroups 获取用户所在的群聊
func (s *GroupService) GetUserGroups(userID uint) (map[string][]vo.GroupChatVO, error) {
	resp := map[string][]vo.GroupChatVO{
		GroupNameMap[Owner]:  make([]vo.GroupChatVO, 0),
		GroupNameMap[Admin]:  make([]vo.GroupChatVO, 0),
		GroupNameMap[Member]: make([]vo.GroupChatVO, 0),
	}

	// 查询用户创建的群组
	var createdGroups []db.Group
	if err := s.db.Where("owner_id = ?", userID).Find(&createdGroups).Error; err != nil {
		return nil, err
	}

	for _, group := range createdGroups {
		resp[GroupNameMap[Owner]] = append(resp[GroupNameMap[Owner]], vo.GroupChatVO{
			ID:          group.ID,
			Name:        group.Name,
			OwnerID:     group.OwnerID,
			Role:        Owner,             // 创建者群主
			GroupAvatar: group.GroupAvatar, // 添加群组头像
			CreatedAt:   group.CreatedAt,
			UpdatedAt:   group.UpdatedAt,
		})
	}

	// 查询用户管理的群组
	var managedGroups []db.GroupMember
	if err := s.db.Where("user_id = ? AND role = ?", userID, Admin).Find(&managedGroups).Error; err == nil {
		for _, member := range managedGroups {
			var group db.Group
			if err := s.db.First(&group, member.GroupID).Error; err == nil {
				resp[GroupNameMap[Admin]] = append(resp[GroupNameMap[Admin]], vo.GroupChatVO{
					ID:          group.ID,
					Name:        group.Name,
					OwnerID:     group.OwnerID,
					Role:        Admin,             // 管理员角色
					GroupAvatar: group.GroupAvatar, // 添加群组头像
					CreatedAt:   group.CreatedAt,
					UpdatedAt:   group.UpdatedAt,
				})
			}
		}
	}

	// 查询用户加入的群组
	var joinedGroups []db.GroupMember
	if err := s.db.Where("user_id = ? AND role = ?", userID, Member).Find(&joinedGroups).Error; err == nil {
		for _, member := range joinedGroups {
			var group db.Group
			if err := s.db.First(&group, member.GroupID).Error; err == nil {
				resp[GroupNameMap[Member]] = append(resp[GroupNameMap[Member]], vo.GroupChatVO{
					ID:          group.ID,
					Name:        group.Name,
					OwnerID:     group.OwnerID,
					Role:        Member,            // 普通成员角色
					GroupAvatar: group.GroupAvatar, // 添加群组头像
					CreatedAt:   group.CreatedAt,
					UpdatedAt:   group.UpdatedAt,
				})
			}
		}
	}

	return resp, nil
}
