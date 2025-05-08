package middle

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"im-system/internal/model/db"
	"time"

	"im-system/internal/config"

	"github.com/dgrijalva/jwt-go"
)

// Claims 自定义的 JWT 载荷
type Claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

func GetRedisUserInfoKey(userId uint) string {
	return fmt.Sprintf("user_login_id:%d", userId)
}

// GetRedisJWTKey 获取 JWT 的 Redis key
func GetRedisJWTKey(userId uint) string {
	return fmt.Sprintf("login:user:jwt:%d", userId)
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

	// 从 Redis 中获取存储的 JWT
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	storedToken, err := config.RedisClient.Get(ctx, GetRedisJWTKey(claims.UserID)).Result()
	if err != nil {
		return 0, errors.New("token已过期或不存在")
	}

	// 比对当前 token 和存储的 token 是否一致
	if storedToken != tokenString {
		return 0, errors.New("token已失效")
	}

	return claims.UserID, nil
}

// SetTokenToRedis 将 token 和用户信息存储到 Redis
func SetTokenToRedis(userID uint, user db.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	redisKey := GetRedisUserInfoKey(userID)
	jwtKey := GetRedisJWTKey(userID)

	marshal, err := json.Marshal(user)
	if err != nil {
		return err
	}

	// 生成 JWT
	token, err := GenerateJWT(userID)
	if err != nil {
		return err
	}

	// 使用 pipeline 批量执行 Redis 命令
	pipe := config.RedisClient.Pipeline()

	// 存储用户基本信息
	pipe.Set(ctx, redisKey, marshal, 72*time.Hour)
	// 存储 JWT token
	pipe.Set(ctx, jwtKey, token, 72*time.Hour)

	// 执行 pipeline
	_, err = pipe.Exec(ctx)
	return err
}
