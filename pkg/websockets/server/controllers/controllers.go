package controllers

import (
	"net/http"
	"time"

	"github.com/azar-intelops/go-websockets-poc/pkg/websockets/server/models"
	"github.com/azar-intelops/go-websockets-poc/pkg/websockets/server/services"
	"github.com/gorilla/websocket"
)

type ChatController struct {
	chatService *services.ChatService
}

func NewChatController(chatService *services.ChatService) *ChatController {
	return &ChatController{
		chatService: chatService,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Add your security checks here to validate the request origin
		// For example, you can check the request's Origin or referer header
		// and return true only if it matches your allowed origins.
		// You can also perform authentication or other security checks.

		return true
	},
}

func (c *ChatController) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	c.chatService.RegisterClient(conn)

	for {
		var message models.Message
		err := conn.ReadJSON(&message)
		if err != nil {
			break
		}

		message.Timestamp = time.Now()

		c.chatService.BroadcastMessage(message)
	}

	c.chatService.UnregisterClient(conn)
}
