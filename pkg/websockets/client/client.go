package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/azar-intelops/go-websockets-poc/pkg/websockets/server/models"
	"github.com/gorilla/websocket"
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
	log.Printf("Connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Error connecting to WebSocket server:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	// Start a goroutine to read messages from the WebSocket connection
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("Error reading message from server:", err)
				return
			}
			log.Printf("Received message: %s", message)
		}
	}()

	// Start a goroutine to send messages to the WebSocket connection
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				// Read message from user input
				var content string
				log.Print("Enter message: ")
				_, err := fmt.Scanln(&content)
				if err != nil {
					log.Println("Error reading user input:", err)
					continue
				}

				// Create a Message struct
				message := models.Message{
					Sender:    "Client",
					Content:   content,
					Timestamp: time.Now(),
				}

				// Send the message as JSON to the WebSocket server
				err = c.WriteJSON(message)
				if err != nil {
					log.Println("Error sending message to server:", err)
					continue
				}
			}
		}
	}()

	// Wait for interrupt signal (e.g., Ctrl+C) to gracefully close the connection
	<-interrupt
	log.Println("Closing WebSocket connection...")
	err = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("Error closing WebSocket connection:", err)
		return
	}
	select {
	case <-done:
	case <-time.After(time.Second):
		log.Println("Timeout waiting for connection to close")
	}
}
