// The teoproxy client-gui example web serve package.
//
// How to build this web application server:
//
//	# Install fyne executible (if not installed)
//	go install fyne.io/fyne/v2/cmd/fyne@latest
//
//	# Go to serve folder inside client-gui folder
//	cd cmd/websocket/client-gui/serve
//
//	# Build web package
//	fyne package -os wasm --sourceDir ../
//
// (or you can use go generate command to build and run this web server)
//
// How to run this web application server:
//
// After build the web server you can start it with next command:
//
//	# Run web server (in development mode)
//	go run .
//
//	# Run web server (in production mode)
//	go run -tags=prod .
//
// How to build executible of this web application server:
//
//	# Build executible
//	go build -tags=prod .
//
//	# Build executible for Linux
//	GOOS=linux go build -tags=prod .
//
//go:generate fyne package -os wasm --sourceDir ../
package main

import (
	"fmt"
	"net/http"

	"github.com/teonet-go/teoproxy/teonet/server"
)

const (
	appShort = "client-gui-serve"
)

func main() {
	// Define a handler function for the HTTP requests
	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	}

	// Register the handler function to handle all requests
	http.HandleFunc("/hello", handler)

	// Create a file server to serve static files from the "wasm" directory
	// fs := http.FileServer(http.Dir("wasm"))
	// http.Handle("/", http.StripPrefix("/", fs))

	// TODO: when embedded filesystem on
	// Register the file server handler to handle all requests
	// http.Handle("/", http.FileServer(http.FS(fs)))
	frontendFS := http.FileServer(http.FS(getFrontendAssets()))
	http.Handle("/", frontendFS)

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
