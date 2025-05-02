package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	ChatType   string             `bson:"chat_type"` // group 或 private
	FromUserID uint               `bson:"from_user_id"`
	ToID       uint               `bson:"to_id"` // 群聊ID或接收者ID
	Content    string             `bson:"content"`
	CreatedAt  time.Time          `bson:"created_at"`
}
