package model

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	FromUserID uint              `bson:"from_user_id"`
	ToUserID   uint              `bson:"to_user_id"`
	Content    string            `bson:"content"`
	CreatedAt  time.Time         `bson:"created_at"`
} 