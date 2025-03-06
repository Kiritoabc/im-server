package vo

import "time"

// MessageVO 表示聊天消息的VO结构体
type MessageVO struct {
	ID         int       `json:"id"`
	Content    string    `json:"content"`
	SenderID   int       `json:"senderId"`
	SenderName string    `json:"senderName"`
	Avatar     string    `json:"avatar"`
	CreatedAt  time.Time `json:"createdAt,omitempty"` // 可以根据实际情况添加时间戳
}

// ChatVO 表示聊天会话的VO结构体
type ChatVO struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	Avatar      string      `json:"avatar"`
	LastMessage string      `json:"lastMessage"`
	Type        string      `json:"type"`
	Messages    []MessageVO `json:"messages"`
}

// ChatsVO 表示所有聊天会话的VO结构体
type ChatsVO []ChatVO
