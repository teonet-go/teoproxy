package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func main() {
	// Start websocket server
	http.HandleFunc("/ws", handleWebSocket)
	log.Fatal(http.ListenAndServe(":8085", nil))

	// Start https websocket server
	// log.Fatal(http.ListenAndServeTLS(":8443", "cert.pem", "key.pem", nil))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}

	// Handle websocket connection
	go handleConnection(conn)
}

func handleConnection(conn *websocket.Conn) {
	defer conn.Close()

	for {
		// Read message from client
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Failed to read message from client:", err)
			break
		}

		// Process message
		processMessage(conn, message)

		// Write response to client
		// err = conn.WriteMessage(websocket.TextMessage, []byte("Message received"))
		// if err != nil {
		// 	log.Println("Failed to write message to client:", err)
		// 	break
		// }
	}
}

func processMessage(conn *websocket.Conn, message []byte) {
	// Process message logic here
	log.Println("Received message:", string(message))

	// Write response to client
	err := conn.WriteMessage(websocket.TextMessage, []byte("Message received"))
	if err != nil {
		log.Println("Failed to write message to client:", err)
	}
}
