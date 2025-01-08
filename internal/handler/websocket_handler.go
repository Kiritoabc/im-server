package handler

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"im-system/internal/service"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketHandler struct {
	clients        map[*websocket.Conn]bool
	broadcast      chan Message
	mu             sync.Mutex
	messageService *service.MessageService
}

type Message struct {
	SenderID   uint   `json:"sender_id"`
	ReceiverID uint   `json:"receiver_id"`
	Content    string `json:"content"`
}

func NewWebSocketHandler(messageService *service.MessageService) *WebSocketHandler {
	return &WebSocketHandler{
		clients:        make(map[*websocket.Conn]bool),
		broadcast:      make(chan Message),
		messageService: messageService,
	}
}

// HandleWebSocket 处理 WebSocket 连接
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	h.mu.Lock()
	h.clients[conn] = true
	h.mu.Unlock()

	for {
		var msg Message
		if err := conn.ReadJSON(&msg); err != nil {
			h.mu.Lock()
			delete(h.clients, conn)
			h.mu.Unlock()
			break
		}

		if err := h.messageService.SaveMessage(msg.SenderID, msg.ReceiverID, msg.Content); err != nil {
			// 处理保存消息的错误
		}
		// TODO: 消息的处理逻辑

		h.broadcast <- msg
	}
}

// StartMessageHandler 启动消息处理
func (h *WebSocketHandler) StartMessageHandler() {
	for {
		msg := <-h.broadcast
		h.mu.Lock()
		for client := range h.clients {
			// todo: 处理消息的发送
			if err := client.WriteJSON(msg); err != nil {
				client.Close()
				delete(h.clients, client)
			}
		}
		h.mu.Unlock()
	}
}
