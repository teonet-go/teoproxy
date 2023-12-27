package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// go:embed wasm/*
// var fs embed.FS

// How to run web application server:
//
//	# Go to main client-gui folder
//	cd client-gui
//
//	# Build web package
//	fyne package -os wasm
//
//	# Run web server
//	go run ./serve
func main() {
	// Define a handler function for the HTTP requests
	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	}

	// Register the handler function to handle all requests
	http.HandleFunc("/hello", handler)

	// Create a file server to serve static files from the "wasm" directory
	fs := http.FileServer(http.Dir("wasm"))

	// Register the file server handler to handle requests starting with "/static/"
	http.Handle("/", http.StripPrefix("/", fs))

	// TODO: when embedded filesystem on
	// Register the file server handler to handle all requests
	//http.Handle("/", http.FileServer(http.FS(fs)))

	// Register websocket server
	http.HandleFunc("/ws", handleWebSocket)

	// Start the web server and listen on port 8080
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}

// func ws() {
// 	// Start websocket server
// 	http.HandleFunc("/ws", handleWebSocket)
// 	// log.Fatal(http.ListenAndServe(":8085", nil))

// 	// Start https websocket server
// 	// log.Fatal(http.ListenAndServeTLS(":8443", "cert.pem", "key.pem", nil))
// }

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
