// Copyright 2023-2024 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The selected code implements a WebSocket server in Go.
//
// The main() function starts the WebSocket server by calling http.HandleFunc()
// to handle "/ws" requests with the handleWebSocket() function. It then starts
// the HTTP server on port 8082.
//
// The handleWebSocket() function upgrades the HTTP request to a WebSocket
// connection using the websocket.Upgrader. It then starts a goroutine to handle
// the WebSocket connection.
//
// The handleConnection() function handles communication with a client over the
// WebSocket connection. It runs in a loop reading messages from the client with
// conn.ReadMessage(), processing the messages by calling processMessage(), and
// closing the connection when done.
//
// The processMessage() function takes the WebSocket connection and the message
// received from the client. It prints the message to the console, writes a
// response "Message received" back to the client over the WebSocket, and
// returns any error.
//
// Overall, this implements a basic WebSocket server that can accept WebSocket
// connections, receive messages, process them, and send responses back over
// the WebSocket. The server runs continuously to handle multiple connections
// and messages over time.
package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// main starts the websocket server.
func main() {
	// Start websocket server
	http.HandleFunc("/ws", handleWebSocket)
	log.Fatal(http.ListenAndServe(":8082", nil))
}

// handleWebSocket handles the WebSocket connection.
//
// It takes in a http.ResponseWriter and a *http.Request as parameters.
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

// handleConnection handles the connection with a client.
//
// It takes a pointer to a websocket.Conn as a parameter.
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
	}
}

// processMessage processes the incoming message received from the client.
//
// It takes in a websocket connection and a byte array representing the message.
// The function prints the received message to the console and writes a response
// message to the client. It returns an error if there was a failure in writing
// the message to the client.
func processMessage(conn *websocket.Conn, message []byte) {
	// Print message to console
	log.Println("Received message:", string(message))

	// Write response to client
	err := conn.WriteMessage(websocket.TextMessage, []byte("Message received"))
	if err != nil {
		log.Println("Failed to write message to client:", err)
	}
}
