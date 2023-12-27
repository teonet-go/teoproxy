package main

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

func main() {
	// WebSocket server URL
	u := url.URL{Scheme: "ws", Host: "localhost:8085", Path: "/ws"}

	// Establish a WebSocket connection
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Error connecting to WebSocket server:", err)
	}
	defer conn.Close()

	// Start a goroutine to receive messages from the server
	go receiveMessages(conn)

	// Send a message to the server
	err = conn.WriteMessage(websocket.TextMessage, []byte("Hello, server!"))
	if err != nil {
		log.Println("Error sending message to WebSocket server:", err)
	}

	// Keep the client running
	select {}
}

func receiveMessages(conn *websocket.Conn) {
	for {
		// Read a message from the server
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error receiving message from WebSocket server:", err)
			return
		}

		// Process the received message
		log.Println("Received message from server:", string(message))
	}
}