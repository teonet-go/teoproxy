package server

import (
	"encoding/base64"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type WsServer struct {
	processMessage []func(conn *websocket.Conn, message []byte)
}

func New(processMessage ...func(conn *websocket.Conn, message []byte)) *WsServer {
	return &WsServer{processMessage: processMessage}
}

func (s *WsServer) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}

	// Handle websocket connection
	go s.handleConnection(conn)
}

func (s *WsServer) handleConnection(conn *websocket.Conn) {
	defer conn.Close()

	for {
		// Read message from client
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Failed to read message from client:", err)
			break
		}

		// Process message
		if len(s.processMessage) == 0 {
			processMessage(conn, message)
			continue
		}
		for _, f := range s.processMessage {
			f(conn, message)
		}
	}
}

func processMessage(conn *websocket.Conn, message []byte) {
	// Process message logic here
	log.Println("Received message:", message, string(message))

	// Write response to client
	sendMessage(conn, []byte("Message received"))
}

// SendMessage sends a message to the websocket client
func sendMessage(conn *websocket.Conn, message []byte) (err error) {
	if err = conn.WriteMessage(websocket.TextMessage,
		[]byte(base64.StdEncoding.EncodeToString(message))); err != nil {
		log.Println("Failed to write message to client:", err)
		return
	}
	log.Println("Message sent to client:", message)
	return
}
