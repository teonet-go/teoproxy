// Copyright 2023-2024 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !wasm

// Package client provides the client implementations for connecting to the
// websocket server.
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

// NewWsClient creates a new instance of the WsClient struct.
//
// The function accepts a variadic parameter `processMessage` of type 
// `func(message []byte) bool`. It returns a pointer to a WsClient instance.
func NewWsClient(processMessage ...func(message []byte) bool) *WsClient {
	return &WsClient{}
}

// SendMessage sends a message to the websocket server.
func (ws *WsClient) SendMessage(message []byte) {
	ws.Conn.WriteMessage(websocket.TextMessage, message)
}

// receiveMessages receives messages from the server.
//
// It reads a message from the WebSocket server and processes it.
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

// Connect opens a WebSocket connection to the server.
// It initializes a WebSocket URL, dials the server, sets the connection
// on the WsClient, starts a goroutine to receive messages,
// sends an initial message, and returns any error.
func (ws *WsClient) Connect() (err error) {
	// WebSocket server URL
	u := url.URL{Scheme: "ws", Host: "localhost:8081", Path: "/ws"}

	// Establish a WebSocket connection
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Println("Error connecting to WebSocket server:", err)
		return
	}
	ws.Conn = conn

	// Start a goroutine to receive messages from the server
	go ws.receiveMessages()

	// Send a message to the server
	err = conn.WriteMessage(websocket.TextMessage, []byte("Hello, server! (inside go)"))
	if err != nil {
		log.Println("Error sending message to WebSocket server:", err)
		return
	}

	return
}
