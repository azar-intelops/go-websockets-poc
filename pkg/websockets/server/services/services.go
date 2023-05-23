package services

import (
	"log"
	"sync"

	"github.com/azar-intelops/go-websockets-poc/pkg/websockets/server/daos"
	"github.com/azar-intelops/go-websockets-poc/pkg/websockets/server/models"
	"github.com/gorilla/websocket"
)

type ChatService struct {
	clients     map[*websocket.Conn]bool
	clientsLock sync.Mutex
	messageDAO  *daos.MessageDAO
}

func NewChatService(messageDAO *daos.MessageDAO) *ChatService {
	return &ChatService{
		clients:    make(map[*websocket.Conn]bool),
		messageDAO: messageDAO,
	}
}

func (cs *ChatService) RegisterClient(client *websocket.Conn) {
	cs.clientsLock.Lock()
	defer cs.clientsLock.Unlock()

	cs.clients[client] = true

	messages := cs.messageDAO.GetMessages()
	for _, message := range messages {
		err := client.WriteJSON(message)
		if err != nil {
			log.Println("Error sending message to client:", err)
		}
	}
}

func (cs *ChatService) UnregisterClient(client *websocket.Conn) {
	cs.clientsLock.Lock()
	defer cs.clientsLock.Unlock()

	delete(cs.clients, client)
}

func (cs *ChatService) BroadcastMessage(message models.Message) {
	cs.clientsLock.Lock()
	defer cs.clientsLock.Unlock()

	for client := range cs.clients {
		err := client.WriteJSON(message)
		if err != nil {
			log.Println("Error sending message to client:", err)
		}
	}

	cs.messageDAO.SaveMessage(message)
}
