package daos

import (
	"sync"

	"github.com/azar-intelops/go-websockets-poc/pkg/websockets/server/models"
)

type MessageDAO struct {
	messages []models.Message
	mutex    sync.Mutex
}

func NewMessageDAO() *MessageDAO {
	return &MessageDAO{}
}

func (dao *MessageDAO) SaveMessage(message models.Message) {
	dao.mutex.Lock()
	defer dao.mutex.Unlock()

	dao.messages = append(dao.messages, message)
}

func (dao *MessageDAO) GetMessages() []models.Message {
	dao.mutex.Lock()
	defer dao.mutex.Unlock()

	return dao.messages
}
