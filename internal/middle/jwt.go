package middle

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"im-system/internal/model/db"
	"time"

	"github.com/dgrijalva/jwt-go"
	"im-system/internal/config"
)

// Claims 自定义的 JWT 载荷
type Claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

func GetRedisUserInfoKey(userId uint) string {
	return fmt.Sprintf("user_login_id:%d", userId)
}

// GenerateJWT 生成 JWT
func GenerateJWT(userID uint) (string, error) {
	claims := Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(), // 72小时后过期
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.JWTSecret))
}

// ValidateJWT 验证 JWT
func ValidateJWT(tokenString string) (uint, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return 0, errors.New("无效的token")
	}

	return claims.UserID, nil
}

// SetTokenToRedis 将 token 和用户信息存储到 Redis
func SetTokenToRedis(userID uint, user db.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	redisKey := GetRedisUserInfoKey(userID)

	marshal, err := json.Marshal(user)
	if err != nil {
		return err
	}

	// 存储用户基本信息
	return config.RedisClient.Set(ctx, redisKey, marshal, 72*time.Hour).Err()
}
