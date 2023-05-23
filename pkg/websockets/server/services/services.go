package services

import (
	"log"
	"sync"

	"github.com/azar-intelops/go-websockets-poc/pkg/websockets/server/daos"
	"github.com/azar-intelops/go-websockets-poc/pkg/websockets/server/models"
	"github.com/gorilla/websocket"
)

type ChatService struct {
	messageDAO *daos.MessageDAO
	clients    map[*websocket.Conn]bool
	mutex      sync.Mutex
}

func NewChatService(messageDAO *daos.MessageDAO) *ChatService {
	return &ChatService{
		messageDAO: messageDAO,
		clients:    make(map[*websocket.Conn]bool),
	}
}

func (s *ChatService) RegisterClient(client *websocket.Conn) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.clients[client] = true
}

func (s *ChatService) UnregisterClient(client *websocket.Conn) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.clients, client)
}

func (s *ChatService) BroadcastMessage(message models.Message) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for client := range s.clients {
		err := client.WriteJSON(message)
		if err != nil {
			log.Println("Error sending message to client:", err)
		}
	}
}
