package main

import (
	"fmt"
	"net/http"

	"github.com/teonet-go/teoproxy/teonet/server"
)

const (
	appShort = "client-gui-serve"
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
	serve, err := server.New(appShort)
	if err != nil {
		fmt.Println("Create teonet proxy server error:", err)
		return
	}
	http.HandleFunc("/ws", serve.HandleWebSocket)

	// Start the web server and listen on port 8080
	err = http.ListenAndServe(":8081", nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}
