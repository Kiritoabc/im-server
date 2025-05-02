package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"im-system/internal/model"
	"im-system/internal/module/kimi"
	"im-system/internal/service"

	"github.com/gin-gonic/gin"
)

type SummaryRequest struct {
	ChatType string `json:"chat_type" binding:"required"` // group 或 private
	UserID   uint   `json:"user_id" binding:"required"`   // 当前用户ID
	ToID     uint   `json:"to_id" binding:"required"`     // 群聊ID或聊天对象ID
}

type SummaryResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Summary string `json:"summary"`
	} `json:"data"`
}

type ChatSummaryHandler struct {
	messageService *service.MessageService
}

func NewChatSummaryHandler(messageService *service.MessageService) *ChatSummaryHandler {
	return &ChatSummaryHandler{
		messageService: messageService,
	}
}

// HandleChatSummary 处理聊天总结请求
func (h *ChatSummaryHandler) HandleChatSummary(c *gin.Context) {
	var req SummaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的请求参数",
		})
		return
	}

	// 验证请求参数
	if err := validateRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	// 获取聊天记录
	messages, err := h.getChatMessages(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取聊天记录失败",
		})
		return
	}

	// 如果没有聊天记录
	if len(messages) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "success",
			"data": gin.H{
				"summary": "暂无聊天记录",
			},
		})
		return
	}

	// 构建上下文内容
	context := buildChatContext(messages)

	// 调用 Kimi API 进行总结
	summary, err := h.callKimiAPI(c.Request.Context(), context)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "生成总结失败",
		})
		return
	}

	// 返回响应
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"summary": summary,
		},
	})
}

func validateRequest(req *SummaryRequest) error {
	if req.ChatType != "group" && req.ChatType != "private" {
		return errors.New("无效的聊天类型")
	}
	if req.UserID == 0 {
		return errors.New("无效的用户ID")
	}
	if req.ToID == 0 {
		return errors.New("无效的目标ID")
	}
	return nil
}

func (h *ChatSummaryHandler) getChatMessages(req *SummaryRequest) ([]model.Message, error) {
	// 使用 messageService 获取消息记录
	messages, err := h.messageService.GetMessageWithSenderInfo(req.ChatType, req.UserID, req.ToID)
	if err != nil {
		return nil, err
	}

	// 转换为统一的消息格式
	var result []model.Message
	for _, msg := range messages {
		message := model.Message{
			FromUserID: msg.SenderID,
			Content:    msg.Content,
			CreatedAt:  msg.CreatedAt,
			ChatType:   req.ChatType,
			ToID:       req.ToID,
		}
		result = append(result, message)
	}

	return result, nil
}

func buildChatContext(messages []model.Message) string {
	var context string
	for _, msg := range messages {
		context += fmt.Sprintf("用户%d: %s\n", msg.FromUserID, msg.Content)
	}
	return context
}

func (h *ChatSummaryHandler) callKimiAPI(ctx context.Context, con string) (string, error) {
	// 创建 ChatCompletionStream 请求
	stream, err := kimi.KimiClient.CreateChatCompletionStream(ctx, &kimi.ChatCompletionStreamRequest{
		Messages: []*kimi.Message{
			{
				Role: kimi.RoleSystem,
				Content: &kimi.Content{Text: "你是一个专业的对话总结助手。请对以下对话内容进行简要总结，" +
					"突出重点，使用简洁的语言。总结要包含：1. 主要话题 2. 关键信息 3. 重要结论或决定"},
			},
			{
				Role: kimi.RoleUser,
				// {"id":1,"senderId":2,
				//"senderName":"小B",
				//"avatar":"https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQ3w2fqb71MsCj97IKLAUXoI6BS4IfeCeEoq_XGS3X2CErGlYyP4xxX4eQ&s",
				//"content":"你好呀","messageType":"private","groupId":null,"receiverId":1}
				Content: &kimi.Content{Text: "请总结以下对话内容,用户的对话在content中，类似json，里面的content是对话的内容，senderName是用户的姓名，请帮我总结一下这些用户聊了些什么：\n" + con},
			},
		},
		Model:       kimi.ModelMoonshot8K,
		MaxTokens:   4096,
		N:           1,
		Temperature: "0.3",
	})

	if err != nil {
		log.Printf("创建聊天流失败: %v", err)
		return "", err
	}
	defer stream.Close()

	var fullSummary string

	// 读取 Moonshot API 的响应
	for chunk := range stream.C {
		deltaContent := chunk.GetDeltaContent()
		if deltaContent != "" {
			fullSummary += deltaContent
		}
	}

	// 检查流式读取是否出错
	if err := stream.Err(); err != nil {
		log.Printf("流处理错误: %v", err)
		return "", err
	}

	return fullSummary, nil
}
