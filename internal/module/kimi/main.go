//go:build !generate_models_file && !context_cache
// +build !generate_models_file,!context_cache

package kimi

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"im-system/internal/config"
	"log"
	"net/http"
	"time"
)

var (
	_ Logger           = moonshot{}
	_ CustomHTTPClient = moonshot{}
)

type moonshot struct {
	baseUrl string
	key     string
	client  *http.Client
	log     func(ctx context.Context, caller string, request *http.Request, response *http.Response, elapse time.Duration)
}

func (m moonshot) BaseUrl() string      { return m.baseUrl }
func (m moonshot) Key() string          { return m.key }
func (m moonshot) Client() *http.Client { return m.client }

func (m moonshot) Log(ctx context.Context, caller string, request *http.Request, response *http.Response, elapse time.Duration) {
	m.log(ctx, caller, request, response, elapse)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// AiClient 全局变量
var client = NewClient[moonshot](moonshot{
	baseUrl: "https://api.moonshot.cn/v1",
	key:     "sk", // 替换为你的 Moonshot API 密钥
	client:  http.DefaultClient,
	log: func(ctx context.Context, caller string, request *http.Request, response *http.Response, elapse time.Duration) {
		log.Printf("[%s] %s %s", caller, request.URL, elapse)
	},
})

func HandleWebSocket(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	defer ws.Close()

	ctx := context.Background()

	// todo: 假设从请求中获取用户 ID 或会话 ID
	userID := "user_123" // 这里可以根据实际情况从请求中获取用户 ID 或生成会话 ID

	// 读取前端发送的消息
	_, message, err := ws.ReadMessage()
	if err != nil {
		log.Println("WebSocket read failed:", err)
		return
	}

	// 获取当前时间戳
	timestamp := time.Now().UnixNano() // 使用纳秒级时间戳作为分数

	// 将用户提问保存到 Redis Sorted Set 中
	userMessage := string(message)
	err = config.RedisClient.ZAdd(ctx, "user:"+userID+":history", &redis.Z{
		Score:  float64(timestamp),    // 时间戳作为分数
		Member: "user:" + userMessage, // 用户消息
	}).Err()
	if err != nil {
		log.Println("Failed to save user question to Redis:", err)
	}

	// 创建 ChatCompletionStream 请求
	stream, err := client.CreateChatCompletionStream(ctx, &ChatCompletionStreamRequest{
		Messages: []*Message{
			{
				Role:    RoleSystem,
				Content: &Content{Text: "你是 Kimi，由 Moonshot AI 提供的人工智能助手，你更擅长中文和英文的对话。你会为用户提供安全，有帮助，准确的回答。同时，你会拒绝一切涉及恐怖主义，种族歧视，黄色暴力等问题的回答。Moonshot AI 为专有名词，不可翻译成其他语言。"},
			},
			{
				Role:    RoleUser,
				Content: &Content{Text: userMessage},
			},
		},
		Model:       ModelMoonshot8K,
		MaxTokens:   4096,
		N:           1,
		Temperature: "0.3",
	})

	if err != nil {
		log.Println("CreateChatCompletionStream failed:", err)
		return
	}
	defer stream.Close()

	var fullResponse string

	// 流式读取 Moonshot API 的响应并发送给前端
	for chunk := range stream.C {
		deltaContent := chunk.GetDeltaContent()
		if deltaContent != "" {
			// 通过 WebSocket 发送流式数据
			if err := ws.WriteMessage(websocket.TextMessage, []byte(deltaContent)); err != nil {
				log.Println("WebSocket write failed:", err)
				return
			}
			// 拼接完整的返回信息
			fullResponse += deltaContent
		}
	}

	// 获取当前时间戳
	timestamp = time.Now().UnixNano()

	// 将 AI 的完整返回信息保存到 Redis Sorted Set 中
	err = config.RedisClient.ZAdd(ctx, "user:"+userID+":history", &redis.Z{
		Score:  float64(timestamp),          // 时间戳作为分数
		Member: "assistant:" + fullResponse, // AI 返回信息
	}).Err()
	if err != nil {
		log.Println("Failed to save assistant response to Redis:", err)
	}

	// 检查流式读取是否出错
	if err := stream.Err(); err != nil {
		log.Println("Stream error:", err)
		return
	}
}
