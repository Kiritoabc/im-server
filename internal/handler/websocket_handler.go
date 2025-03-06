package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	"im-system/internal/service"
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
	broadcast      chan Message
	mu             sync.Mutex
	messageService *service.MessageService
}

type Message struct {
	Id         int    `json:"id"`
	SenderId   int    `json:"senderId"`
	ReceiverID int    `json:"receiverId"`
	SenderName string `json:"senderName"`
	Avatar     string `json:"avatar"`
	Content    string `json:"content"`
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
	handleWebSocket(ctx.Writer, ctx.Request)
}

// handleWebSocket 处理WebSocket连接
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
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
		receiver := strconv.FormatInt(int64(msg.ReceiverID), 10)
		// 发送私聊消息给目标用户,不打印也就是说没找到
		log.Println(receiver)
		if client, OK := clients[receiver]; OK {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Println("Write message error:", err)
				client.Close()
				delete(clients, receiver)
			}
			log.Println(msg)
		}
	}
}
