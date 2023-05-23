package main

import (
	"log"
	"net/http"

	"github.com/azar-intelops/go-websockets-poc/pkg/websockets/server/controllers"
	"github.com/azar-intelops/go-websockets-poc/pkg/websockets/server/daos"
	"github.com/azar-intelops/go-websockets-poc/pkg/websockets/server/services"
)

func main() {
	messageDAO := daos.NewMessageDAO()
	chatService := services.NewChatService(messageDAO)
	chatController := controllers.NewChatController(chatService)

	http.HandleFunc("/ws", chatController.WebSocketHandler)

	log.Println("Starting chat application...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
