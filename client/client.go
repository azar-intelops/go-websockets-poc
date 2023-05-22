package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gorilla/websocket"
)

var (
	username string
	conn     *websocket.Conn
)

// Message represents a message sent to the server
type Message struct {
	Username string `json:"username"`
	Content  string `json:"content"`
}

func main() {
	// Get username from command-line input
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter username: ")
	username, _ = reader.ReadString('\n')
	username = strings.TrimSpace(username)

	// Connect to the WebSocket server
	connectToServer()

	// Handle incoming messages from the server
	go handleIncomingMessages()

	// Send user input to the server
	handleUserInput()

	// Wait for SIGINT or SIGTERM to gracefully close the connection
	waitForExitSignal()
}

func connectToServer() {
	var err error
	conn, _, err = websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handleIncomingMessages() {
	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Error while reading message:", err)
			break
		}
		fmt.Printf("[%s]: %s\n", msg.Username, msg.Content)
	}
}

func handleUserInput() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		scanner.Scan()
		message := scanner.Text()
		if strings.ToLower(message) == "/quit" {
			break
		}

		// Send user input as a message to the server
		msg := Message{
			Username: username,
			Content:  message,
		}
		err := conn.WriteJSON(msg)
		if err != nil {
			log.Println("Error while sending message:", err)
		}
	}
}

func waitForExitSignal() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	<-interrupt
	fmt.Println("Closing connection...")
	conn.Close()
	fmt.Println("Disconnected.")
}
