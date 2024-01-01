//go:build !wasm

package client

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

// WsClient is javascript websocket client to use in wasm application.
type WsClient struct {
	*websocket.Conn
}

func NewWsClient(processMessage ...func(message []byte)) *WsClient {
	return &WsClient{}
}

func (ws *WsClient) SendMessage(message []byte) {
	// Send a message to the websocket server
	ws.Conn.WriteMessage(websocket.TextMessage, message)
}

func (ws *WsClient) receiveMessages() {
	for {
		// Read a message from the server
		_, message, err := ws.Conn.ReadMessage()
		if err != nil {
			log.Println("Error receiving message from WebSocket server:", err)
			return
		}

		// Process the received message
		log.Println("Received message from server:", string(message))
	}
}

// Connect to websocket server in native application
func (ws *WsClient) Connect() (err error) {
	// WebSocket server URL
	u := url.URL{Scheme: "ws", Host: "localhost:8081", Path: "/ws"}

	// Establish a WebSocket connection
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Println("Error connecting to WebSocket server:", err)
		return
	}
	// defer conn.Close()
	ws.Conn = conn

	// Start a goroutine to receive messages from the server
	go ws.receiveMessages()

	// Send a message to the server
	err = conn.WriteMessage(websocket.TextMessage, []byte("Hello, server! (inside go)"))
	if err != nil {
		log.Println("Error sending message to WebSocket server:", err)
		return
	}

	// Keep the client running
	// select {}

	return
}
