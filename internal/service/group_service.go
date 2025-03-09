package service

import (
	"errors"
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

const (
	defaultUserAvatar = "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQ3w2fqb71MsCj97IKLAUXoI6BS4IfeCeEoq_XGS3X2CErGlYyP4xxX4eQ&s"
	GroupTypePublic   = "public"
	GroupTypePrivate  = "private"
)

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

// GetMyAllGroups 获取用户的所有群聊
func (s *GroupService) GetMyAllGroups(userId uint) (vo.GroupChatList, error) {
	// 查询所有和我有关的群聊
	var groupMembers []db.GroupMember
	if err := s.db.Where("user_id =?", userId).Find(&groupMembers).Error; err != nil {
		return vo.GroupChatList{}, err
	}
	// 构建响应数据
	var groupChatList vo.GroupChatList
	var createGroups []vo.GroupChatDetail
	var managedGroups []vo.GroupChatDetail
	var joinedGroups []vo.GroupChatDetail
	for _, groupMember := range groupMembers {
		// 查询群的基本信息
		var group db.Group
		if err := s.db.First(&group, groupMember.GroupID).Error; err != nil {
			return vo.GroupChatList{}, err
		}
		// todo: 查询群的统计信息
		// 查询群成员信息
		var members []db.GroupMember
		if err := s.db.Where("group_id =?", groupMember.GroupID).Find(&members).Error; err != nil {
			return vo.GroupChatList{}, err
		}
		var groupStatus vo.GroupStats
		groupStatus.Active = len(members)
		groupStatus.Male = 1
		groupStatus.Local = 0
		var previewMembers []vo.PreviewMember
		for _, member := range members {
			previewMembers = append(previewMembers, vo.PreviewMember{
				ID:   member.UserID,
				Name: member.Nickname,
				// todo: 查询用户头像
				Avatar: defaultUserAvatar,
				Role:   member.Role,
			})
		}
		// 判断是否是创建者，管理员，普通成员
		if groupMember.Role == Owner {
			createGroups = append(createGroups, vo.GroupChatDetail{
				ID:          group.ID,
				Name:        group.Name,
				Avatar:      defaultUserAvatar,
				MemberCount: len(members),
				// todo: 添加群聊统计信息
				Category:       "游戏交友",
				Announcement:   "暂无",
				Description:    "暂无",
				Stats:          groupStatus,
				PreviewMembers: previewMembers,
			})
		} else if groupMember.Role == Admin {
			managedGroups = append(managedGroups, vo.GroupChatDetail{
				ID:          group.ID,
				Name:        group.Name,
				Avatar:      defaultUserAvatar,
				MemberCount: len(members),
				// todo: 添加群聊统计信息
				Category:       "游戏交友",
				Announcement:   "暂无",
				Description:    "暂无",
				Stats:          groupStatus,
				PreviewMembers: previewMembers,
			})
		} else if groupMember.Role == Member {
			joinedGroups = append(joinedGroups, vo.GroupChatDetail{
				ID:          group.ID,
				Name:        group.Name,
				Avatar:      defaultUserAvatar,
				MemberCount: len(members),
				// todo: 添加群聊统计信息
				Category:       "游戏交友",
				Announcement:   "暂无",
				Description:    "暂无",
				Stats:          groupStatus,
				PreviewMembers: previewMembers,
			})
		}
	}
	groupChatList.CreatedGroups = createGroups
	groupChatList.ManagedGroups = managedGroups
	groupChatList.JoinedGroups = joinedGroups
	return groupChatList, nil
}

// GetGroupMembers 获取群聊成员
func (s *GroupService) GetGroupMembers(groupId string) ([]vo.UserVO, error) {
	// 查询群聊成员
	var groupMembers []db.GroupMember
	if err := s.db.Where("group_id =?", groupId).Find(&groupMembers).Error; err != nil {
		return nil, err
	}
	// 构建响应数据
	var userVOs []vo.UserVO
	for _, groupMember := range groupMembers {
		// 查询用户信息
		var user db.User
		if err := s.db.First(&user, groupMember.UserID).Error; err != nil {
			return nil, err
		}
		userVOs = append(userVOs, vo.UserVO{
			ID:          user.ID,
			Username:    user.Username,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
			AvatarURL:   user.AvatarURL,
			Bio:         user.Bio,
			Gender:      user.Gender,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			City:        user.City,
		})
	}
	return userVOs, nil
}

// InviteGroup 邀请好友加入群聊
func (s *GroupService) InviteGroup(userId uint, groupId uint, friendIds []uint) error {
	// 查询群聊信息
	var group db.Group
	if err := s.db.First(&group, groupId).Error; err != nil {
		return err
	}
	// 判断用户是否是群主或管理员
	var groupMember db.GroupMember
	if err := s.db.Where("group_id = ? AND user_id =?", groupId, userId).First(&groupMember).Error; err != nil {
		return err
	}
	if groupMember.Role != Owner && groupMember.Role != Admin {
		return errors.New("只有群主或管理员可以邀请好友加入群聊")
	}

	// 判断好友是否已经在群聊中,如果在群聊中，不拉入该好友，否则拉入群聊
	for _, friendId := range friendIds {
		var member db.GroupMember
		if err := s.db.Where("group_id =? AND user_id =?", groupId, friendId).First(&member).Error; err == nil {
			continue
		}
		// 查询好友的基本信息
		var friend db.User
		if err := s.db.First(&friend, friendId).Error; err != nil {
			return err
		}
		// 好友不在群聊中，加入群聊
		member = db.GroupMember{
			GroupID:  groupId,
			UserID:   friendId,
			Nickname: friend.Username,
			Role:     Member,
			Level:    1,
		}
		if err := s.db.Create(&member).Error; err != nil {
			return err
		}
	}
	return nil
}
