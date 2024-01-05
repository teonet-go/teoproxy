// Copyright 2023-2024 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The selected code implements a WebSocket client in Go.
//
// It starts by defining the WebSocket server URL to connect to. It uses the
// gorilla/websocket package to establish a connection to the server.
//
// It then starts a goroutine to receive messages from the server asynchronously.
// This allows the client to both send and receive messages independently.
//
// The client sends an initial "Hello, server!" message to the server after
// connecting. It uses the websocket.Conn WriteMessage method for this.
//
// The receiveMessages goroutine loops continuously to read incoming messages
// from the server using the ReadMessage method. Each received message is
// printed to the log.
//
// So in summary, this code sets up a WebSocket client, connects to a server,
// sends an initial message, and then processes responses continuously as they
// arrive from the server. The goroutine allows bi-directional communication 
// with asynchronous receive handling.
//
// The main purpose is to demonstrate a simple WebSocket client in Go that can
// send and receive messages from a server over a persistent WebSocket
// connection. It uses goroutines and channels to handle the receive loop
// concurrently.
package main

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

func main() {
	// WebSocket server URL
	u := url.URL{Scheme: "ws", Host: "localhost:8082", Path: "/ws"}

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
