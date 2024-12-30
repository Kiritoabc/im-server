package middle

import (
	"context"
	"encoding/json"
	"im-system/internal/model/db"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"im-system/internal/config"
)

// AuthMiddleware JWT 验证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 忽略注册和登录请求
		if c.Request.URL.Path == "/im-server/register" || c.Request.URL.Path == "/im-server/login" {
			c.Next()
			return
		}

		// 从请求头中获取 token
		token := c.Request.Header.Get("token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供token"})
			c.Abort()
			return
		}

		// 验证 token
		userID, err := ValidateJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的token"})
			c.Abort()
			return
		}

		// 从 Redis 中获取用户信息
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		cacheUserInfo, err := config.RedisClient.Get(ctx, getRedisUserInfoKey(userID)).Result()
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户信息已过期，请重新登录"})
			c.Abort()
			return
		}
		// todo: 将userInfo 反序列化
		userInfo := &db.User{}
		if err = json.Unmarshal([]byte(cacheUserInfo), &userInfo); err != nil {
			config.Logger.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
			return
		}

		// 将用户ID和用户信息存储到上下文中
		c.Set("user_id", userID)
		c.Set("user_info", userInfo)

		// 继续处理请求
		c.Next()
	}
}
