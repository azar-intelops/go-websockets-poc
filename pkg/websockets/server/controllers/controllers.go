package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/azar-intelops/go-websockets-poc/pkg/websockets/server/models"
	"github.com/azar-intelops/go-websockets-poc/pkg/websockets/server/services"
	"github.com/gorilla/websocket"
)

type ChatController struct {
	chatService *services.ChatService
}

// ChatService handles chat-related operations


func NewChatController(chatService *services.ChatService) *ChatController {
	return &ChatController{
		chatService: chatService,
	}
}

func (c *ChatController) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// Add your security checks here to validate the request origin
			// For example, you can check the request's Origin or referer header
			// and return true only if it matches your allowed origins.
			// You can also perform authentication or other security checks.

			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	c.chatService.RegisterClient(conn)

	for {
		var message models.Message
		err := conn.ReadJSON(&message)
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		message.Timestamp = time.Now()

		c.chatService.BroadcastMessage(message)
	}

	c.chatService.UnregisterClient(conn)
}
