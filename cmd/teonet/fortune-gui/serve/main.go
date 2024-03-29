// Copyright 2023-2024 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The teoproxy fortune-gui example web serve package.
//
// How to build this web application server:
//
//	# Install fyne executible (if not installed)
//	go install fyne.io/fyne/v2/cmd/fyne@latest
//
//	# Go to serve folder inside fortune-gui folder
//	cd cmd/websocket/fortune-gui/serve
//
//	# Build web package
//	fyne package -os wasm --appVersion=0.0.3 --sourceDir ../
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
//go:generate fyne package -os wasm --appVersion=0.0.3 --sourceDir ../
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/NYTimes/gziphandler"
	"github.com/teonet-go/teoproxy/teonet/server"
	"golang.org/x/crypto/acme/autocert"
)

const (
	appShort   = "fortune-gui-serve"
	appName    = "Fortune-gui web server"
	appVersion = "0.0.3"
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
	var monitor, laddr string
	var gzip bool
	//
	flag.StringVar(&domain, "domain", "", "domain name to process HTTP/s server")
	flag.StringVar(&laddr, "laddr", "localhost:8081", "local address of http, used if domain doesn't set")
	flag.StringVar(&monitor, "monitor", "", "teonet monitor address")
	flag.BoolVar(&gzip, "gzip", false, "gzip http files")
	flag.Parse()

	// Define Hello handler function for the HTTP requests
	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	}

	// Register the Hello handler function to handle all requests
	http.HandleFunc("/hello", handler)

	// Create a file server to serve static files from the "wasm" directory
	var frontendFS http.Handler
	if gzip {
		frontendFS = gziphandler.GzipHandler(http.FileServer(http.FS(getFrontendAssets())))
	} else {
		frontendFS = http.FileServer(http.FS(getFrontendAssets()))
	}
	http.Handle("/", frontendFS)

	// Register teonet proxy server handler
	serve, err := server.New(appShort, &server.TeonetMonitor{
		Addr:       monitor,
		AppName:    appName,
		AppShort:   appShort,
		AppVersion: appVersion,
	})
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
