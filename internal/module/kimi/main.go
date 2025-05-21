//go:build !generate_models_file && !context_cache
// +build !generate_models_file,!context_cache

package kimi

import (
	"context"
	"encoding/json"
	"im-system/internal/config"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
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

// KimiClient 全局变量
var KimiClient = NewClient[moonshot](moonshot{
	baseUrl: "https://api.moonshot.cn/v1",
	key:     "", // 替换为你的 Moonshot API 密钥
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

	// 定义用户消息结构
	type UserMessage struct {
		UserId  json.Number `json:"userId"` // 使用 json.Number 可以接受字符串或数字类型
		Id      int         `json:"id"`
		Role    string      `json:"role"`
		Content string      `json:"content"`
	}

	// 读取前端发送的消息
	_, message, err := ws.ReadMessage()
	if err != nil {
		log.Println("WebSocket read failed:", err)
		return
	}

	// 解析 JSON 消息
	var userMsg UserMessage
	if err := json.Unmarshal(message, &userMsg); err != nil {
		log.Println("Failed to parse user message:", err)
		return
	}

	// 获取当前时间戳
	timestamp := time.Now().UnixNano() // 使用纳秒级时间戳作为分数

	// 将用户提问保存到 Redis Sorted Set 中
	userMessage := userMsg.Content
	err = config.RedisClient.ZAdd(ctx, "user:"+userID+":history", &redis.Z{
		Score:  float64(timestamp),    // 时间戳作为分数
		Member: "user:" + userMessage, // 用户消息
	}).Err()
	if err != nil {
		log.Println("Failed to save user question to Redis:", err)
	}

	// 创建非流式 ChatCompletion 请求
	completion, err := KimiClient.CreateChatCompletion(ctx, &ChatCompletionRequest{
		Messages: []*Message{
			{
				Role:    RoleSystem,
				Content: &Content{Text: "你是 Kimi，由 Moonshot AI 提供的人工智能助手，你更擅长中文和英文的对话。你会为用户提供安全，有帮助，准确的回答。同时，你会拒绝一切涉及恐怖主义，种族歧视，黄色暴力等问题的回答。Moonshot AI 为专有名词，不可翻译成其他语言。请使用 Markdown 格式来组织你的回答，以提供更好的阅读体验。"},
			},
			{
				Role:    RoleUser,
				Content: &Content{Text: userMessage},
			},
		},
		Model:          ModelMoonshot8K,
		MaxTokens:      4096,
		N:              1,
		Temperature:    "0.3",
		ResponseFormat: ResponseFormatText,
	})

	if err != nil {
		log.Println("CreateChatCompletion failed:", err)
		return
	}

	// 获取完整响应
	fullResponse := completion.Choices[0].Message.Content.Text

	// 处理Markdown格式
	// 1. 确保段落之间有空行
	fullResponse = strings.ReplaceAll(fullResponse, "\n\n\n", "\n\n") // 先统一多个连续换行为两个
	fullResponse = strings.ReplaceAll(fullResponse, "\n", "\n\n")     // 确保每个换行后有空行

	// 2. 处理列表项
	lines := strings.Split(fullResponse, "\n")
	for i := 0; i < len(lines); i++ {
		if strings.HasPrefix(strings.TrimSpace(lines[i]), "- ") || strings.HasPrefix(strings.TrimSpace(lines[i]), "* ") {
			if i > 0 && !strings.HasPrefix(strings.TrimSpace(lines[i-1]), "- ") && !strings.HasPrefix(strings.TrimSpace(lines[i-1]), "* ") {
				lines[i] = "\n" + lines[i]
			}
		}
	}
	fullResponse = strings.Join(lines, "\n")

	// 3. 处理代码块
	fullResponse = strings.ReplaceAll(fullResponse, "```", "\n```\n")

	// 4. 处理标题
	for i := 6; i >= 1; i-- {
		prefix := strings.Repeat("#", i) + " "
		fullResponse = strings.ReplaceAll(fullResponse, "\n"+prefix, "\n\n"+prefix)
	}

	// 通过 WebSocket 一次性发送完整响应
	if err := ws.WriteMessage(websocket.TextMessage, []byte(fullResponse)); err != nil {
		log.Println("WebSocket write failed:", err)
		return
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
}
