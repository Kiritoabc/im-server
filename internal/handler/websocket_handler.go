package handler

import (
	"encoding/json"
	"im-system/internal/config"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"

	"im-system/internal/service"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 存储所有连接的客户端
var clients = make(map[string]*websocket.Conn)

type WebSocketHandler struct {
	clients        map[string]*websocket.Conn
	groupClients   map[int][]*websocket.Conn
	broadcast      chan Message
	mu             sync.Mutex
	messageService *service.MessageService
}

type Message struct {
	Id          int    `json:"id"`
	SenderId    int    `json:"senderId"`
	ReceiverID  int    `json:"receiverId"` // 用于私聊
	GroupID     int    `json:"groupId"`    // 用于群聊
	SenderName  string `json:"senderName"`
	Avatar      string `json:"avatar"`
	Content     string `json:"content"`
	MessageType string `json:"messageType"` // "private" 或 "group"
	CreatedAt   string `json:"createdAt"`   // 2025-05-07T16:17:21+08:00
}

func NewWebSocketHandler(messageService *service.MessageService) *WebSocketHandler {
	return &WebSocketHandler{
		clients:        clients,
		broadcast:      make(chan Message),
		messageService: messageService,
	}
}

// SendMessage 发送消息
func (h *WebSocketHandler) SendMessage(ctx *gin.Context) {
	h.handleWebSocket(ctx.Writer, ctx.Request)
}

// handleWebSocket 处理WebSocket连接
func (h *WebSocketHandler) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	// 读取客户端发送的用户名
	_, username, err := conn.ReadMessage()
	if err != nil {
		log.Println("Read username error:", err)
		return
	}
	// 判断是否重复
	clients[string(username)] = conn

	for {
		// 读取客户端发送的消息
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read message error:", err)
			delete(clients, string(username))
			break
		}
		// 解析消息
		var msg Message
		err = json.Unmarshal(message, &msg)
		if err != nil {
			log.Println("Unmarshal message error:", err)
			continue
		}
		// 保存消息到数据库
		// 判断是群消息还是私聊消息
		if msg.MessageType == "group" {
			// 处理群聊消息
			if err := h.messageService.SaveGroupMessage(uint(msg.SenderId), uint(msg.GroupID), string(message)); err != nil {
				// 处理保存失败的情况
				log.Println("保存消息失败", err)
				config.Logger.Errorf("保存消息失败: %v", err)
			}
		}
		if msg.MessageType == "private" {
			// 处理私聊消息
			if err := h.messageService.SaveMessage(uint(msg.SenderId), uint(msg.ReceiverID), string(message)); err != nil {
				// 处理保存失败的情况
				log.Println("保存消息失败", err)
				config.Logger.Errorf("保存消息失败: %v", err)
			}
		}

		log.Println(msg)
		if msg.MessageType == "private" {
			// 处理私聊消息
			receiver := strconv.FormatInt(int64(msg.ReceiverID), 10)
			if client, OK := clients[receiver]; OK {
				err := client.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Println("Write message error:", err)
					client.Close()
					delete(clients, receiver)
				}
			}
		} else if msg.MessageType == "group" {
			// 处理群聊消息
			// 使用 messageService 获取群组成员
			members, err := h.messageService.GetGroupMembers(msg.GroupID)
			if err != nil {
				log.Println("Error fetching group members:", err)
				continue
			}
			// 发送消息给群组中的所有成员
			for _, member := range members {
				receiver := strconv.FormatInt(int64(member.UserID), 10)

				if client, OK := clients[receiver]; OK && int(member.UserID) != msg.SenderId {
					err := client.WriteMessage(websocket.TextMessage, message)
					if err != nil {
						log.Println("Write message error:", err)
						client.Close()
					}
				}
			}
		}
	}
}
