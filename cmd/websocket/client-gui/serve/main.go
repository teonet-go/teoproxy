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
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/teonet-go/teoproxy/teonet/server"
	"golang.org/x/crypto/acme/autocert"
	"github.com/NYTimes/gziphandler"
)

const (
	appShort = "client-gui-serve"
)

var domain string

// main is the entry point of the program.
//
// It parses application parameters, defines a handler function for HTTP requests,
// registers the handler function to handle all requests, creates a file server
// to serve static files, registers a websocket server, starts an HTTPS server
// if a domain is set, or starts an HTTP server if a domain is not set.
func main() {

	// Parse application parameters
	var laddr string
	//
	flag.StringVar(&domain, "domain", "", "domain name to process HTTP/s server")
	flag.StringVar(&laddr, "laddr", "localhost:8081", "local address of http, used if domain doesn't set")
	flag.Parse()

	// Define Hello handler function for the HTTP requests
	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	}

	// Register the Hello handler function to handle all requests
	http.HandleFunc("/hello", handler)

	// Create a file server to serve static files from the "wasm" directory
	frontendFS := gziphandler.GzipHandler(http.FileServer(http.FS(getFrontendAssets())))
	http.Handle("/", frontendFS)

	// Register websocket server
	serve, err := server.New(appShort)
	if err != nil {
		fmt.Println("Create teonet proxy server error:", err)
		return
	}
	http.HandleFunc("/ws", serve.HandleWebSocket)

	// Start HTTPS server if domain is set
	if len(domain) > 0 {

		// Redirect HTTP requests to HTTPS
		go func() {
			err := http.ListenAndServe(":80", http.HandlerFunc(redirectTLS))
			if err != nil {
				log.Fatalf("ListenAndServe error: %v", err)
			}
		}()

		// Start HTTPS server and create certificate for domain
		log.Println("Start https serve with domain:", domain)
		log.Fatal(http.Serve(autocert.NewListener(domain), nil))
		return
	}

	// Start HTTP server
	log.Println("Start http serve at:", laddr)
	log.Fatalln(http.ListenAndServe(laddr, nil))
}

// redirectTLS redirects the HTTP request to HTTPS.
//
// It takes in the http.ResponseWriter and http.Request as parameters.
func redirectTLS(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+domain+":443"+r.RequestURI,
		http.StatusMovedPermanently)
}
