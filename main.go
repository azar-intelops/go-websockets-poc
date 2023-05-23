package main

import (
	"log"
	"net/http"

	"github.com/azar-intelops/go-websockets-poc/pkg/websockets/server/controllers"
	"github.com/azar-intelops/go-websockets-poc/pkg/websockets/server/daos"
	"github.com/azar-intelops/go-websockets-poc/pkg/websockets/server/services"
)

func main() {
	// Initialize DAO
	messageDAO := daos.NewMessageDAO()

	// Initialize Service
	chatService := services.NewChatService(messageDAO)

	// Initialize Controller
	chatController := controllers.NewChatController(chatService)

	// Define routes and handlers
	http.HandleFunc("/ws", chatController.WebSocketHandler)

	log.Println("Starting chat application...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
