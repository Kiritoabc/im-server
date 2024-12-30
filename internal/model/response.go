package model

import (
	"github.com/gin-gonic/gin"
)

// Response 通用响应结构体
type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`  // 可选数据字段
	Error   string      `json:"error,omitempty"` // 可选错误字段
}

// Success 创建成功响应
func Success(message string, data interface{}) Response {
	return Response{
		Message: message,
		Data:    data,
	}
}

// Error 创建错误响应
func Error(message string) Response {
	return Response{
		Error: message,
	}
}

// SendResponse 发送响应
func SendResponse(c *gin.Context, status int, response Response) {
	c.JSON(status, response)
}
