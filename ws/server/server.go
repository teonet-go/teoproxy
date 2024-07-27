// Copyright 2023-2024 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package server provides the WebSocket server implementation.
package server

import (
	"encoding/base64"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// WsServer is a WebSocket server that handles WebSocket connections.
// It contains a processMessage field which is a slice of functions to process
// incoming WebSocket messages.
type WsServer struct {
	processMessage []func(conn *websocket.Conn, message []byte)
	onClose        func(conn *websocket.Conn)
}

// New creates a new WsServer instance with the provided message processing
// functions. The processMessage functions will be called to handle each
// incoming WebSocket message.
func New(onClose func(conn *websocket.Conn), processMessage ...func(conn *websocket.Conn, message []byte)) *WsServer {
	return &WsServer{processMessage: processMessage, onClose: onClose}
}

// HandleWebSocket handles websocket requests by upgrading
// the HTTP connection to a WebSocket connection.
func (s *WsServer) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		// Allow cross-origin websocket connections
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}

	// Handle websocket connection
	go s.handleConnection(conn)
}

// handleConnection handles an incoming WebSocket connection.
// It reads messages from the client, processes them by calling functions in
// processMessage, and runs until the connection is closed.
func (s *WsServer) handleConnection(conn *websocket.Conn) {
	defer conn.Close()

	log.Printf("ws client connected %p %v", conn, conn.RemoteAddr())
	for {
		// Read message from client
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("failed to read message from client:", err)
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

	log.Printf("ws client disconnected %p, %s", conn, conn.RemoteAddr())
	if s.onClose != nil {
		s.onClose(conn)
	}
}

// processMessage handles incoming WebSocket messages from clients.
// It logs the message, processes it, and writes a response.
func processMessage(conn *websocket.Conn, message []byte) {
	// Print message to console
	log.Println("received message:", message, string(message))

	// Write response to client
	sendMessage(conn, []byte("Message received"))
}

// SendMessage sends a message to the websocket client
// sendMessage sends a message to the websocket client.
// It encodes the message as base64 text and writes it to the client.
// Returns any error from writing the message.
func sendMessage(conn *websocket.Conn, message []byte) (err error) {
	if err = conn.WriteMessage(websocket.TextMessage,
		[]byte(base64.StdEncoding.EncodeToString(message))); err != nil {
		log.Println("failed to write message to client:", err)
		return
	}
	log.Println("message sent to client:", message)
	return
}
