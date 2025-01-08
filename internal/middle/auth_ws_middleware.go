package middle

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// AuthWSMiddleware WebSocket 鉴权中间件
func AuthWSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取 Token
		token := c.Query("token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供token"})
			c.Abort()
			return
		}

		// 验证 Token
		userID, err := ValidateJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的token"})
			c.Abort()
			return
		}

		// 将用户ID存储到上下文中
		c.Set("user_id", userID)

		// 继续处理请求
		c.Next()
	}
}
