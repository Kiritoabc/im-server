package service

import (
	"errors"
	"gorm.io/gorm"
	"im-system/internal/model/db"
	"im-system/internal/model/vo"
)

type NotificationService struct {
	db *gorm.DB
}

func NewNotificationService() *NotificationService {
	return &NotificationService{
		db: db.DB,
	}
}

// GetNotifications 获取用户的通知
func (s *NotificationService) GetNotifications(userID uint, notificationType string) ([]vo.NotificationVO, error) {
	var notifications []db.Notification
	query := s.db.Where("receiver_id = ?", userID)

	// 根据类型过滤通知
	if notificationType != "" {
		query = query.Where("type = ?", notificationType)
	}

	if err := query.Find(&notifications).Error; err != nil {
		return nil, err
	}

	// 创建 NotificationVO 列表
	var notificationVOs []vo.NotificationVO
	for _, notification := range notifications {
		var sender db.User
		// 查询发送添加好友请求人的信息
		if err := s.db.First(&sender, notification.SenderID).Error; err == nil {
			notificationVOs = append(notificationVOs, vo.NotificationVO{
				ID:         notification.ID,
				UserID:     notification.SenderID,
				ReceiverID: notification.ReceiverID,
				Type:       notification.Type,
				Content:    notification.Content,
				IsRead:     notification.IsRead,
				CreatedAt:  notification.CreatedAt,
				Sender:     sender, // 发送请求的用户信息
			})
		}
	}

	return notificationVOs, nil
}

// AcceptFriendRequest 接受好友请求
func (s *NotificationService) AcceptFriendRequest(notificationID uint, groupId uint) error {
	var notification db.Notification
	// 查询通知
	if err := s.db.First(&notification, notificationID).Error; err != nil {
		return errors.New("通知不存在")
	}
	// 检查通知是否已读
	if notification.IsRead {
		return nil
	}

	// 更新通知状态为已读
	notification.IsRead = true
	if err := s.db.Save(&notification).Error; err != nil {
		return err
	}

	// 查询请求者的基本信息
	var requester db.User
	if err := s.db.First(&requester, notification.SenderID).Error; err != nil {
		return errors.New("请求者信息不存在")
	}

	// 创建请求者的好友关系
	friendship := db.Friendship{
		UserID:   notification.ReceiverID, // 被请求者的用户ID
		FriendID: notification.SenderID,   // 请求者的用户ID
		Status:   "accepted",              // 状态为已接受
		Remark:   requester.Username,      // 使用请求者的用户名作为备注
		GroupID:  groupId,                 // 分组ID
	}

	if err := s.db.Create(&friendship).Error; err != nil {
		return err
	}

	return nil
}

// RejectFriendRequest 拒绝好友请求
func (s *NotificationService) RejectFriendRequest(notificationID uint) error {
	var notification db.Notification
	if err := s.db.First(&notification, notificationID).Error; err != nil {
		return errors.New("通知不存在")
	}

	// 更新通知状态为已读
	notification.IsRead = true
	notification.Type = "rejected"
	if err := s.db.Save(&notification).Error; err != nil {
		return err
	}

	// 这里可以添加其他拒绝请求的逻辑，将发送者创建的friendship删除
	if err := s.db.Model(&db.Friendship{}).
		Where("user_id =? and friend_id =?", notification.SenderID, notification.ReceiverID).
		Delete(&db.Friendship{}).Error; err != nil {
		return err
	}

	return nil
}

// GetSentNotifications 获取用户发出的所有通知请求
func (s *NotificationService) GetSentNotifications(userID uint) ([]vo.NotificationVO, error) {
	var notifications []db.Notification
	if err := s.db.Where("sender_id = ?", userID).Find(&notifications).Error; err != nil {
		return nil, err
	}

	// 创建 NotificationVO 列表
	var notificationVOs []vo.NotificationVO
	for _, notification := range notifications {
		var receiver db.User
		// 查询接收者的信息
		if err := s.db.First(&receiver, notification.ReceiverID).Error; err == nil {
			notificationVOs = append(notificationVOs, vo.NotificationVO{
				ID:         notification.ID,
				UserID:     notification.SenderID,
				ReceiverID: notification.ReceiverID,
				Type:       notification.Type,
				Content:    notification.Content,
				IsRead:     notification.IsRead,
				CreatedAt:  notification.CreatedAt,
				Sender:     receiver, // 接收者的信息
			})
		}
	}

	return notificationVOs, nil
}
